package git

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"k8s.io/klog"
)

const (
	defaultRetries = 3
)

// New returns a git Interface
func New(path string) (Interface, error) {
	g, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}
	return &gitClient{
		git:  g,
		path: path,
	}, nil
}

type gitClient struct {
	git     string
	path    string
	repo    string
	repoURL *url.URL
}

// WithRepo implements git.Interface
func (c *gitClient) WithRepo(repo string) Interface {
	c.repo = repo
	c.repoURL = nil
	return c
}

// Clone implements git.Interface
// TODO(liubog2008): support cache
// TODO(liubog2008): support credential
func (c *gitClient) Clone() error {
	localRepo, err := c.localRepo()
	if err != nil {
		return err
	}

	if _, err := os.Stat(localRepo); os.IsNotExist(err) {
		if err := os.MkdirAll(localRepo, os.ModePerm); err != nil {
			return err
		}

		output, err := retry(defaultRetries, "", c.git, "clone", c.repo, localRepo)
		if err != nil {
			return fmt.Errorf("git clone err: %v, output: %s", err, string(output))
		}
	} else {
		output, err := retry(defaultRetries, localRepo, c.git, "fetch")
		if err != nil {
			return fmt.Errorf("git fetch error: %v. output: %s", err, string(output))
		}
	}

	return nil
}

// Checkout implements git.Interface
func (c *gitClient) Checkout(commit string) error {
	localRepo, err := c.localRepo()
	if err != nil {
		return err
	}
	if _, err := retry(1, localRepo, c.git, "checkout", commit); err != nil {
		return err
	}
	return nil
}

func (c *gitClient) localRepo() (string, error) {
	if c.repoURL == nil {
		u, err := parse(c.repo)
		if err != nil {
			return "", err
		}
		c.repoURL = u
	}
	localRepo := filepath.Join(c.path, c.repoURL.Host, c.repoURL.Path)
	return localRepo, nil
}

func parse(repo string) (*url.URL, error) {
	if repo == "" {
		return nil, fmt.Errorf("repo path is empty")
	}
	// NOTE(liubog2008): Now SCP-like URL is not supported
	return url.Parse(repo)
}

func retry(retries int, dir, cmd string, arg ...string) ([]byte, error) {
	var (
		lastOutput []byte
		lastError  error
	)
	sleepTime := time.Second
	for i := 0; i < retries; i++ {
		c := exec.Command(cmd, arg...)
		c.Dir = dir
		b, err := c.CombinedOutput()
		if err == nil {
			return b, nil
		}
		klog.Warningf("Running %s %v returned error %v with output %s.", cmd, arg, err, string(b))
		time.Sleep(sleepTime)
		sleepTime *= 2
		lastOutput = b
		lastError = err
	}
	return lastOutput, lastError
}
