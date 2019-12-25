package mario

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	marioFileCotent = []byte(`mario test file`)
)

func newFakeGitRepo() (string, error) {
	d, err := ioutil.TempDir("", "git")
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(filepath.Join(d, v1alpha1.MarioFile), marioFileCotent, 0644); err != nil {
		return "", err
	}

	g, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}

	initCmd := exec.Command(g, "init")
	initCmd.Dir = d
	if err := initCmd.Run(); err != nil {
		return "", err
	}

	addCmd := exec.Command(g, "add", ".")
	addCmd.Dir = d
	if err := addCmd.Run(); err != nil {
		return "", err
	}

	commitCmd := exec.Command(g, "commit", "-m", "test commit")
	commitCmd.Dir = d
	if err := commitCmd.Run(); err != nil {
		return "", err
	}

	return d, nil
}

const (
	// TODO(liubog2008): change to use a random string
	testToken = "asdfvve"
)

func TestServe(t *testing.T) {
	gitRepo, err := newFakeGitRepo()
	require.NoError(t, err, "can't create fake git repo")
	defer os.RemoveAll(gitRepo)

	localPath, err := ioutil.TempDir("", "git")
	require.NoError(t, err, "can't new temp repo")
	defer os.RemoveAll(localPath)

	waitCh := make(chan struct{})

	m, err := newWithWait(localPath, gitRepo, ":8080", testToken, waitCh)
	require.NoError(t, err, "can't create mario interface")

	require.NoError(t, m.Checkout("master"), "can't fetch and checkout repo")

	stopCh := make(chan struct{})

	go func() {
		require.NoError(t, m.Serve(stopCh), "can't serve")
	}()

	<-waitCh

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	require.NoError(t, err, "new GET request err")

	req.Header.Set(authKey, tokenType+" "+testToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "do GET request err")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "status code should be OK")

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err, "ready body has no error")
	assert.Equal(t, marioFileCotent, body, "cotent should be equal")
}
