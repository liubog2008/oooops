// Package mario deifnes interface to download git repo and serve mario file
package mario

type Interface interface {
	// Verify is used for verifying git status
	Verify()
}
