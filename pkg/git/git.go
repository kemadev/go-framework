// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
	"github.com/go-git/go-git/v6/storage/memory"
)

const (
	BranchMain = "main"
	BranchNext = "next"
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

func (g *Service) GetRemoteGitRepo(remoteURL string) (*git.Repository, error) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: remoteURL,
	})
	if err != nil {
		return nil, fmt.Errorf("error opening git repository: %w", err)
	}

	return repo, nil
}
