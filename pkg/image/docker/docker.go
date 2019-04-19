package docker

import (
	"os"
	"os/exec"
)

// Interface defines docker interface to build and push image
// TODO(liubog2008): add credential interface
type Interface interface {
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
}

func (d *dockerClient) Build(dockerfile, repository string) error {
	cmd := exec.Command(d.docker, "build", "-t", repository, "-f", dockerfile, d.path)
	cmd.Dir = d.path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *dockerClient) Push(repository string) error {
	cmd := exec.Command(d.docker, "push", repository)
	cmd.Dir = d.path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
