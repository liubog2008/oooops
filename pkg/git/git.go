package git

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"k8s.io/klog"
)

const (
	defaultRetries = 3
)

var (
	ErrMissingGitRepo = errors.New("git repo is not set")
)

// New returns a git Interface
func New(path string) (Interface, error) {
	g, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := os.MkdirAll(path, 0744); err != nil {
			return nil, err
		}
	} else {
		if !stat.IsDir() {
			return nil, fmt.Errorf("path for git is not dir")
		}
	}

	return &gitClient{
		git:  g,
		path: path,
	}, nil
}

type gitClient struct {
	git  string
	path string

	remote string
}

// WithRepo implements git.Interface
func (c *gitClient) WithRepo(repo string) error {
	_, err := parse(repo)
	if err != nil {
		return fmt.Errorf("can't parse repo of git: %v", err)
	}

	c.remote = repo

	return nil
}

// Clone implements git.Interface
// TODO(liubog2008): support cache
// TODO(liubog2008): support credential
func (c *gitClient) Fetch(ref string) error {
	isEmpty, err := isEmptyDir(c.path)
	if err != nil {
		return err
	}
	if isEmpty {
		if err := retry(defaultRetries, c.path, c.git, "clone", c.remote, "."); err != nil {
			return fmt.Errorf("git clone error: %v", err)
		}
	}

	args := []string{"fetch"}
	if ref != "" {
		args = append(args, "origin", ref)
	}

	if err := retry(defaultRetries, c.path, c.git, args...); err != nil {
		return fmt.Errorf("git fetch error: %v", err)
	}

	return nil
}

func (c *gitClient) Clean() error {
	if err := retry(1, c.path, c.git, "clean", "-f"); err != nil {
		return err
	}

	return nil
}

// Checkout implements git.Interface
// Checkout will call git clean before checkout
func (c *gitClient) Checkout(commit string) error {
	if err := retry(1, c.path, c.git, "checkout", commit); err != nil {
		return err
	}

	return nil
}

func parse(repo string) (*url.URL, error) {
	if repo == "" {
		return nil, fmt.Errorf("repo path is empty")
	}
	// NOTE(liubog2008): Now SCP-like URL is not supported
	return url.Parse(repo)
}

func retry(retries int, dir, cmd string, args ...string) error {
	var (
		lastError error
	)
	sleepTime := time.Second
	for i := 0; i < retries; i++ {
		klog.Infof("Trying [%s %v] %v times", cmd, strings.Join(args, " "), i)
		c := exec.Command(cmd, args...)
		c.Dir = dir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err == nil {
			return nil
		}
		time.Sleep(sleepTime)
		sleepTime *= 2
		lastError = err
	}
	return lastError
}

func isEmptyDir(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if _, err := f.Readdirnames(1); err == io.EOF {
		return true, nil
	}
	return false, err
}
