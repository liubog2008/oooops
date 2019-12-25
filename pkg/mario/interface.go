// Package mario deifnes interface to download git repo and serve mario file
package mario

// Interface defines these functions of mario:
// - Pull and checkout code
// - Serve v1alpha1.Mario file as http API
type Interface interface {
	Checkout(ref string) error

	Serve(stopCh chan struct{}) error
}
