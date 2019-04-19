package source

import (
	"fmt"

	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	"github.com/liubog2008/oooops/pkg/source/git"
)

type sourceManager struct {
	git git.Interface
}

// New returns a source Manager
func New(path string) (Manager, error) {
	g, err := git.New(path)
	if err != nil {
		return nil, err
	}
	return &sourceManager{
		git: g,
	}, nil
}

func (m *sourceManager) Fetch(s *v1alpha1.CodeSource) error {
	return m.FetchByGit(&s.Git)
}

func (m *sourceManager) FetchByGit(s *v1alpha1.GitSource) error {
	ref := s.Matches
	switch s.Type {
	case v1alpha1.GitBranch:
		if s.Matches == "" {
			ref = "master"
		}
	case v1alpha1.GitRelease:
		if s.Matches != "" {
			return fmt.Errorf("Release tag is not found")
		}
	case v1alpha1.GitRevision:
		if s.Matches == "" {
			ref = "master"
		}
	case v1alpha1.GitPullRequest:
		return fmt.Errorf("It has not been supported yet")
	default:
		return fmt.Errorf("Unknown type %s", s.Type)
	}

	return m.gitCheckout(s.Repository, ref)
}

func (m *sourceManager) gitCheckout(repo, ref string) error {
	g := m.git.WithRepo(repo)
	if err := g.Fetch(ref); err != nil {
		return err
	}
	if err := g.Checkout("origin/" + ref); err != nil {
		return err
	}
	return nil
}
