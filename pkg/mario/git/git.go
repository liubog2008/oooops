// Package git defines git interface and its implement
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/klog"
)

const (
	retryTimes = 1
)

// Interface defines git interface which is used by mario
type Interface interface {
	Verify(remote, ref string) error
}

type gitCmd struct {
	git        string
	workingDir string
}

// New returns a git Interface
func New(dir string) (Interface, error) {
	g, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(filepath.Clean(dir))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := os.MkdirAll(dir, 0744); err != nil {
			return nil, err
		}
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not dir", dir)
	}

	return &gitCmd{
		git:        g,
		workingDir: dir,
	}, nil
}

func (c *gitCmd) Verify(remote, ref string) error {
	isClean, err := c.IsClean()
	if err != nil {
		return err
	}
	if !isClean {
		return fmt.Errorf("working dir is not clean, please checking your git script")
	}
	if err := c.RemoteIsMatched(remote); err != nil {
		return err
	}
	if err := c.CommitIsMatched(ref); err != nil {
		return err
	}
	return nil
}

// IsClean check whether working dir is clean by running
// command "git status --porcelain"
func (c *gitCmd) IsClean() (bool, error) {
	output, err := c.retry(retryTimes, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	if len(output) == 0 {
		return true, nil
	}
	return false, nil
}

func (c *gitCmd) RemoteIsMatched(remote string) error {
	output, err := c.retry(retryTimes, "remote", "get-url", "origin")
	if err != nil {
		return err
	}
	current := strings.TrimSpace(string(output))
	if current == remote {
		return nil
	}
	return fmt.Errorf(
		"remote is not matched, current is [%s], expected is [%s]",
		current,
		remote,
	)
}

func (c *gitCmd) CommitIsMatched(ref string) error {
	// ^{} means dereference the tag recursively until a non-tag object is found
	commitOutput, err := c.retry(retryTimes, "rev-parse", ref+"^{}")
	if err != nil {
		return err
	}
	commit := strings.TrimSpace(string(commitOutput))
	headOutput, err := c.retry(retryTimes, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	head := strings.TrimSpace(string(headOutput))
	if commit == head {
		return nil
	}
	return fmt.Errorf(
		"ref commit is not matched, current is [%s], expected is [%s(%s)]",
		head,
		commit,
		ref,
	)
}

func (c *gitCmd) retry(retries int, args ...string) ([]byte, error) {
	var (
		lastError error
	)
	sleepTime := time.Second
	for i := 0; i < retries; i++ {
		klog.Infof("Trying [%s %v] %v times", c.git, strings.Join(args, " "), i)
		output, err := runCommand(c.workingDir, c.git, args...)
		if err != nil {
			klog.Errorf("Failed to run [%s %v]: %v\n--- git ---\n%s--- git ---", c.git, strings.Join(args, " "), err, output)
			lastError = err
			time.Sleep(sleepTime)
			sleepTime *= 2
			continue
		}
		klog.Infof("run [%s %v] success\n--- git ---\n%s--- git ---", c.git, strings.Join(args, " "), output)
		return output, nil
	}
	return nil, lastError
}

func runCommand(dir, cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	c.Dir = dir
	return c.CombinedOutput()
}
