package git

type Interface interface {
	// WithRepo sets repo of client
	WithRepo(repo string) error
	// Fetch fetches whole commits of the repo
	// If repo doesn't exist, clone first
	// If ref is not empty, just fetch ref
	Fetch(ref string) error

	// Checkout checkout to specifid commit
	Checkout(commit string) error

	Clean() error
}
