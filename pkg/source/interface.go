package source

import "github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"

// Manager defines source manager
type Manager interface {
	// Fetch fetches code from source
	Fetch(s *v1alpha1.CodeSource) error
}
