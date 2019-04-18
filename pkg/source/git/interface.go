package git

// Interface defines git interface to fetch code
type Interface interface {
	WithRepo(repo string) Interface

	Clone() error

	Checkout(commit string) error
}
