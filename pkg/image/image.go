package image

import (
	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	"github.com/liubog2008/oooops/pkg/image/docker"
)

type imageManager struct {
	docker docker.Interface
}

// New returns a image manager
func New(path string) (Manager, error) {
	d, err := docker.New(path)
	if err != nil {
		return nil, err
	}
	return &imageManager{
		docker: d,
	}, nil
}

// BuildAndPush implements image.Manager
func (m *imageManager) BuildAndPush(im *v1alpha1.Image) error {
	d := m.docker.WorkDir(im.WorkDir)
	if err := d.Build(im.Path, im.Repository); err != nil {
		return err
	}
	return d.Push(im.Repository)
}
