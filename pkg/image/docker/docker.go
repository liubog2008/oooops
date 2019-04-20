package docker

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Interface defines docker interface to build and push image
// TODO(liubog2008): add credential interface
type Interface interface {
	// WorkDir defines work dir
	WorkDir(dir string) Interface

	// Build builds a image
	Build(dockerfile, repository string) error

	// Push pushes built image to remote
	Push(repository string) error
}

// New returns a image interface
func New(path string) (Interface, error) {
	d, err := exec.LookPath("docker")
	if err != nil {
		return nil, err
	}
	return &dockerClient{
		docker: d,
		path:   path,
	}, nil

}

type dockerClient struct {
	docker string
	path   string
	dir    string
}

func (d *dockerClient) WorkDir(dir string) Interface {
	nd := *d
	nd.dir = dir
	return &nd
}

func (d *dockerClient) Build(dockerfile, repository string) error {
	cmd := exec.Command(d.docker, "build", "-t", repository, "-f", dockerfile, ".")
	cmd.Dir = filepath.Join(d.path, d.dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *dockerClient) Push(repository string) error {
	cmd := exec.Command(d.docker, "push", repository)
	cmd.Dir = filepath.Join(d.path, d.dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
