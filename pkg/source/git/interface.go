package git

// Interface defines git interface to fetch code
type Interface interface {
	// WithRepo sets repo of client
	WithRepo(repo string) Interface

	// Fetch fetches whole commits of the repo
	// If repo doesn't exist, clone first
	// If ref is not empty, just fetch ref
	Fetch(ref string) error

	// Checkout checkout to specifid commit
	Checkout(commit string) error
}
