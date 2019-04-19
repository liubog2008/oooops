package image

import "github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"

// Manager defines an image manager
type Manager interface {
	// BuildAndPush builds and pushes image to remote
	BuildAndPush(image *v1alpha1.Image) error
}
