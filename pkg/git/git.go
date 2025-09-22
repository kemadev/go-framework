// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package git

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/svu/v3/pkg/svu"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/storage/memory"
)

var (
	ErrRemoteURLNotFound = fmt.Errorf("remote URL not found")
	ErrGitRepoNil        = fmt.Errorf("git repository is nil")
	ErrBranchesNil       = fmt.Errorf("branches are nil")
	ErrCurrBrancheNil    = fmt.Errorf("current branch is nil")
	ErrGitTreeNil        = fmt.Errorf("git tree is nil")
)

type Service struct {
	repo *git.Repository
}

func NewGitService() *Service {
	return &Service{}
}

func NewGitServiceWithRepo(repo *git.Repository) *Service {
	return &Service{repo: repo}
}

func (g *Service) GetGitRepo() (*git.Repository, error) {
	if g.repo != nil {
		return g.repo, nil
	}

	repo, err := git.PlainOpenWithOptions(
		".",
		&git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: false},
	)
	if err != nil {
		return nil, fmt.Errorf("error opening git repository: %w", err)
	}

	g.repo = repo

	return g.repo, nil
}

func (g *Service) GetGitBranches() (*storer.ReferenceIter, error) {
	repo, err := g.GetGitRepo()
	if err != nil {
		return nil, fmt.Errorf("error getting git repository: %w", err)
	}

	if repo == nil {
		return nil, ErrGitRepoNil
	}

	branches, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("error getting branches: %w", err)
	}

	if branches == nil {
		return nil, ErrBranchesNil
	}

	return &branches, nil
}

func (g *Service) GetGitCurrentBranchName() (string, error) {
	repo, err := g.GetGitRepo()
	if err != nil {
		return "", fmt.Errorf("error getting git repository: %w", err)
	}

	if repo == nil {
		return "", ErrGitRepoNil
	}

	currentBranch, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("error getting current branch: %w", err)
	}

	if currentBranch == nil {
		return "", ErrCurrBrancheNil
	}

	return currentBranch.Name().Short(), nil
}

func (g *Service) GetGitHead() (*plumbing.Reference, error) {
	repo, err := g.GetGitRepo()
	if err != nil {
		return nil, fmt.Errorf("error getting git repository: %w", err)
	}

	if repo == nil {
		return nil, ErrGitRepoNil
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("error getting repository head: %w", err)
	}

	if head == nil {
		return nil, ErrCurrBrancheNil
	}

	return head, nil
}

func (g *Service) GetGitHeadTree() (*object.Tree, error) {
	repo, err := g.GetGitRepo()
	if err != nil {
		return nil, fmt.Errorf("error getting git repository: %w", err)
	}

	if repo == nil {
		return nil, ErrGitRepoNil
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("error getting repository head: %w", err)
	}

	if head == nil {
		return nil, ErrCurrBrancheNil
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("error getting repository commit: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("error getting repository tree: %w", err)
	}

	if tree == nil {
		return nil, fmt.Errorf(
			"error getting repository tree: %w",
			ErrGitTreeNil,
		)
	}

	return tree, nil
}

func (g *Service) GetGitBasePath() (string, error) {
	repo, err := g.GetGitRepo()
	if err != nil {
		return "", fmt.Errorf("error getting git repository: %w", err)
	}

	return g.GetGitBasePathWithRepo(repo)
}

func (g *Service) GetGitBasePathWithRepo(repo *git.Repository) (string, error) {
	remote, err := repo.Remote("origin")
	if err != nil {
		return "", fmt.Errorf("error getting remote: %w", err)
	}

	if len(remote.Config().URLs) == 0 {
		return "", ErrRemoteURLNotFound
	}

	basePath := strings.TrimPrefix(remote.Config().URLs[0], "git@")
	basePath = strings.TrimPrefix(basePath, "https://")
	basePath = strings.TrimSuffix(basePath, ".git")

	return basePath, nil
}

func (g *Service) TagSemver() (bool, error) {
	currentVersion, err := svu.Current()
	if err != nil {
		return false, fmt.Errorf("error getting current version: %w", err)
	}

	slog.Debug("got version", slog.String("current-version", currentVersion))

	nextVersion, err := svu.Next()
	if err != nil {
		return false, fmt.Errorf("error getting next version: %w", err)
	}

	slog.Debug("got version", slog.String("next-version", nextVersion))

	if currentVersion == nextVersion {
		return true, nil
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return false, fmt.Errorf("error getting git repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return false, fmt.Errorf("error getting HEAD reference: %w", err)
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

func (g *Service) GetRemoteGitRepo(remoteURL string) (*git.Repository, error) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: remoteURL,
	})
	if err != nil {
		return nil, fmt.Errorf("error opening git repository: %w", err)
	}

	return repo, nil
}
