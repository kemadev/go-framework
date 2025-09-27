// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var ErrMalformedVersion = errors.New("version string is malformed")

const (
	CaptureGroupKeyMajor         string = "major"
	CaptureGroupKeyMinor         string = "minor"
	CaptureGroupKeyPatch         string = "patch"
	CaptureGroupKeyPreRelease    string = "prerelease"
	CaptureGroupKeyPreType       string = "pretype"
	CaptureGroupKeyPreMajor      string = "premajor"
	CaptureGroupKeyPreMinor      string = "preminor"
	CaptureGroupKeyPrePatch      string = "prepatch"
	CaptureGroupKeyBuildMetadata string = "buildmetadata"
)

// RegExp adapted from https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
const SemverRegex string = `^v?(?P<` + CaptureGroupKeyMajor + `>0|[1-9]\d*)\.(?P<` + CaptureGroupKeyMinor + `>0|[1-9]\d*)\.(?P<` + CaptureGroupKeyPatch + `>0|[1-9]\d*)(?:-(?P<` + CaptureGroupKeyPreRelease + `>(?:(?P<` + CaptureGroupKeyPreType + `>[a-zA-Z][0-9a-zA-Z-]*)\.)?(?P<` + CaptureGroupKeyPreMajor + `>0|[1-9]\d*)\.(?P<` + CaptureGroupKeyPreMinor + `>0|[1-9]\d*)\.(?P<` + CaptureGroupKeyPrePatch + `>0|[1-9]\d*)|(?:[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*)))?(?:\+(?P<` + CaptureGroupKeyBuildMetadata + `>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

type PreReleaseType string

var (
	PreReleaseTypeNone  PreReleaseType = ""
	PreReleaseTypeDev   PreReleaseType = "dev"
	PreReleaseTypeNext  PreReleaseType = "next"
	PreReleaseTypeAlpha PreReleaseType = "alpha"
	PreReleaseTypeBeta  PreReleaseType = "beta"
)

type PreRelease struct {
	Type    PreReleaseType
	Version MajMinPatch
}

type MajMinPatch struct {
	Major int
	Minor int
	Patch int
}

type Version struct {
	Version    MajMinPatch
	PreRelease PreRelease
	BuildMeta  string
}

func (v Version) Compare(other Version) int {
	cmp := v.Version.Compare(other.Version)
	if cmp != 0 {
		return cmp
	}

	if v.PreRelease.Type == PreReleaseTypeNone && other.PreRelease.Type == PreReleaseTypeNone {
		return 0
	}

	if v.PreRelease.Type == PreReleaseTypeNone && other.PreRelease.Type != PreReleaseTypeNone {
		return 1
	}

	if v.PreRelease.Type != PreReleaseTypeNone && other.PreRelease.Type == PreReleaseTypeNone {
		return -1
	}

	return v.PreRelease.Compare(other.PreRelease)
}

func (m MajMinPatch) Compare(other MajMinPatch) int {
	if m.Major != other.Major {
		if m.Major < other.Major {
			return -1
		}

		return 1
	}

	if m.Minor != other.Minor {
		if m.Minor < other.Minor {
			return -1
		}

		return 1
	}

	if m.Patch != other.Patch {
		if m.Patch < other.Patch {
			return -1
		}

		return 1
	}

	return 0
}

func (p PreRelease) Compare(other PreRelease) int {
	if p.Type != other.Type {
		if p.Type < other.Type {
			return -1
		}

		return 1
	}

	return p.Version.Compare(other.Version)
}

func (v Version) LessThan(other Version) bool {
	return v.Compare(other) < 0
}

func (v Version) LessThanOrEqual(other Version) bool {
	return v.Compare(other) <= 0
}

func (v Version) GreaterThan(other Version) bool {
	return v.Compare(other) > 0
}

func (v Version) GreaterThanOrEqual(other Version) bool {
	return v.Compare(other) >= 0
}

func (v Version) Equal(other Version) bool {
	return v.Compare(other) == 0
}

func (p PreReleaseType) String() string {
	return string(p)
}

func (v Version) String() string {
	if v.PreRelease.Type == PreReleaseTypeNone {
		return fmt.Sprintf("%d.%d.%d", v.Version.Major, v.Version.Minor, v.Version.Patch)
	}

	return fmt.Sprintf("%d.%d.%d-%s.%d.%d.%d",
		v.Version.Major, v.Version.Minor, v.Version.Patch,
		v.PreRelease.Type.String(),
		v.PreRelease.Version.Major, v.PreRelease.Version.Minor, v.PreRelease.Version.Patch)
}

func Parse(str string) (Version, error) {
	reg, err := regexp.Compile(SemverRegex)
	if err != nil {
		return Version{}, fmt.Errorf("error compiling regex %q: %w", SemverRegex, err)
	}

	matches := reg.FindStringSubmatch(str)
	if matches == nil {
		return Version{}, fmt.Errorf("version %q: %w", str, ErrMalformedVersion)
	}

	names := reg.SubexpNames()
	result := make(map[string]string)

	for i, name := range names {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	major, err := strconv.Atoi(result[CaptureGroupKeyMajor])
	if err != nil {
		return Version{}, fmt.Errorf("version %q: invalid major version: %w", str, err)
	}

	minor, err := strconv.Atoi(result[CaptureGroupKeyMinor])
	if err != nil {
		return Version{}, fmt.Errorf("version %q: invalid minor version: %w", str, err)
	}

	patch, err := strconv.Atoi(result[CaptureGroupKeyPatch])
	if err != nil {
		return Version{}, fmt.Errorf("version %q: invalid patch version: %w", str, err)
	}

	var preMajor, preMinor, prePatch int

	if result[CaptureGroupKeyPreType] != "" {
		ver, err := strconv.Atoi(result[CaptureGroupKeyPreMajor])
		if err != nil {
			return Version{}, fmt.Errorf(
				"version %q: invalid prerelease major version: %w",
				str,
				err,
			)
		}

		preMajor = ver

		ver, err = strconv.Atoi(result[CaptureGroupKeyPreMinor])
		if err != nil {
			return Version{}, fmt.Errorf(
				"version %q: invalid prerelease minor version: %w",
				str,
				err,
			)
		}

		preMinor = ver

		ver, err = strconv.Atoi(result[CaptureGroupKeyPrePatch])
		if err != nil {
			return Version{}, fmt.Errorf(
				"version %q: invalid prerelease patch version: %w",
				str,
				err,
			)
		}

		prePatch = ver
	}

	return Version{
		Version: MajMinPatch{
			Major: major,
			Minor: minor,
			Patch: patch,
		},
		PreRelease: PreRelease{
			Type: PreReleaseType(result[CaptureGroupKeyPreType]),
			Version: MajMinPatch{
				Major: preMajor,
				Minor: preMinor,
				Patch: prePatch,
			},
		},
		BuildMeta: result[CaptureGroupKeyBuildMetadata],
	}, nil
}
