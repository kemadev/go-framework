package git

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/kemadev/go-framework/pkg/semver"
)

var ErrBranchNameInvalid = errors.New("branche name not in valid set")

func NextTag(
	branchName string,
	repo *git.Repository,
) (semver.Version, error) {
	currentTag, err := CurrentVersion(repo)
	if err != nil {
		return semver.Version{}, fmt.Errorf("branch name %q: %w", branchName, err)
	}
	nextTag := semver.Version{}

	switch branchName {
	case BranchMain:
		ver, err := NextVersion(repo, currentTag)
		if err != nil {
			return semver.Version{}, fmt.Errorf("branch name %q: %w", branchName, err)
		}

		nextTag = ver
	case BranchNext:
		ver, err := NextPreVersion(repo, currentTag)
		if err != nil {
			return semver.Version{}, fmt.Errorf("branch name %q: %w", branchName, err)
		}

		nextTag = ver
	default:
		return semver.Version{}, fmt.Errorf("branch name %q: %w", branchName, ErrBranchNameInvalid)
	}

	return nextTag, nil
}

func CurrentVersion(repo *git.Repository) (semver.Version, error) {
	currentVersion := semver.Version{}

	tags, err := repo.Tags()

	err = tags.ForEach(func(r *plumbing.Reference) error {
		tagString := r.Name().String()

		ver, err := semver.Parse(tagString)
		if err != nil {
			slog.Warn("ignoring invalid tag", slog.String("tag", tagString))
			return nil
		}

		if ver.GreaterThanOrEqual(currentVersion) {
			currentVersion = ver
		}

		return nil
	})
	if err != nil {
		return semver.Version{}, fmt.Errorf("error parsing tags: %w", err)
	}

	return currentVersion, nil
}

// func NextVersion(repo *git.Repository, current semver.Version) (semver.Version, error) {
// 	currentVersion := semver.Version{}

// 	commits, err := repo.CommitObjects()

// 	err = commits.ForEach(func(o *object.Commit) error {
// 		conventionalcommit.

// 		ver, err := semver.Parse(tagString)
// 		if err != nil {
// 			slog.Warn("ignoring invalid tag", slog.String("tag", tagString))
// 			return nil
// 		}

// 		if ver.GreaterThanOrEqual(currentVersion) {
// 			currentVersion = ver
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return semver.Version{}, fmt.Errorf("error parsing tags: %w", err)
// 	}

// 	return currentVersion, nil
// }

// func NextPreVersion(repo *git.Repository, current semver.Version) (semver.Version, error) {
// 	return current, nil
// }
