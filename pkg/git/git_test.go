package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	gitPath = "github.com/liubog2008/oooops"
)

// TODO(liubog2008): add fake repo
func TestFetch(t *testing.T) {
	d, err := ioutil.TempDir("", "git")
	require.NoError(t, err, "can't create temp dir")
	defer os.RemoveAll(d)

	g, err := New(d)
	require.NoError(t, err, "can't create git interface")

	repo := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), gitPath)
	require.NoError(t, g.WithRepo(repo), "can't set repo")

	require.NoError(t, g.Fetch(""), "can't fetch repo")

	thisfile := filepath.Join(d, "pkg/git/git_test.go")
	_, err = os.Stat(thisfile)
	assert.NoError(t, err, "%s should exist", thisfile)
}
