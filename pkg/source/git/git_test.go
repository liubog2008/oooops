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

func TestClone(t *testing.T) {
	d, err := ioutil.TempDir("", "git")
	require.NoError(t, err, "can't create temp dir")
	g, err := New(d)
	require.NoError(t, err, "can't create git interface")
	repo := fmt.Sprintf("file://%s/src/github.com/liubog2008/oooops", os.Getenv("GOPATH"))
	require.NoError(t, g.WithRepo(repo).Clone(), "can't clone repo")
	thisfile := filepath.Join(d, "github.com/liubog2008/oooops", "pkg/source/git/git_test.go")
	_, err = os.Stat(thisfile)
	assert.NoError(t, err, "%s should exist", thisfile)
}
