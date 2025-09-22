package git

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/svu/v3/pkg/svu"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

var ErrBranchNameInvalid = errors.New("branche name not in valid set")

func (g *Service) TagSemver() (bool, error) {
	repo, err := g.GetGitRepo()
	if err != nil {
		return false, fmt.Errorf("error getting git repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return false, fmt.Errorf("error getting HEAD reference: %w", err)
	}

	currentVersion, err := svu.Current()
	if err != nil {
		return false, fmt.Errorf("error getting current version: %w", err)
	}

	slog.Debug("got version", slog.String("current-version", currentVersion))

	branchName := head.Name().String()

	nextVersion, err := NextTag(branchName, repo)
	if err != nil {
		return false, fmt.Errorf("error getting next version: %w", err)
	}

	slog.Debug("got version", slog.String("next-version", nextVersion))

	if currentVersion == nextVersion {
		return true, nil
	}

	ref, err := repo.CreateTag(nextVersion, head.Hash(), nil)
	if err != nil {
		return false, fmt.Errorf("error creating tag: %w", err)
	}

	slog.Info("tag created", slog.String("tag", ref.Name().Short()))

	return false, nil
}

func (g *Service) PushTag() error {
	repo, err := g.GetGitRepo()
	if err != nil {
		return fmt.Errorf("error getting git repository: %w", err)
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		FollowTags: true,
		Auth: &http.BasicAuth{
			Username: "git",
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})
	if err != nil {
		return fmt.Errorf("error pushing tags: %w", err)
	}

	slog.Debug("pushed tag")

	return nil
}

func NextTag(
	branchName string,
	repo *git.Repository,
) (string, error) {
	switch branchName {
	case BranchMain:
		ver, err := svu.Next()
		if err != nil {
			return "", fmt.Errorf("branch name %q: %w", branchName, err)
		}

		return ver, nil
	case BranchNext:
		ver, err := svu.PreRelease()
		if err != nil {
			return "", fmt.Errorf("branch name %q: %w", branchName, err)
		}

		return ver, nil
	}

	return "", fmt.Errorf("branch name %q: %w", branchName, ErrBranchNameInvalid)
}
