package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrMalformedVersion = errors.New("version string is malformed")

type PreReleaseType int

const (
	PreReleaseTypeNone PreReleaseType = iota
	PreReleaseTypeAlpha
	PreReleaseTypeBeta
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

func (p PreReleaseType) String() string {
	switch p {
	case PreReleaseTypeNone:
		return ""
	case PreReleaseTypeAlpha:
		return "alpha"
	case PreReleaseTypeBeta:
		return "beta"
	default:
		return "unknown"
	}
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

func parseMajMinPatch(majStr string, minStr string, patchStr string) (int, int, int, error) {
	maj, err := strconv.ParseInt(majStr, 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("maj %q: %w", maj, ErrMalformedVersion)
	}

	min, err := strconv.ParseInt(minStr, 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("min %q: %w", min, ErrMalformedVersion)
	}

	patch, err := strconv.ParseInt(patchStr, 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("patch %q: %w", patch, ErrMalformedVersion)
	}

	return int(maj), int(min), int(patch), nil
}

func parseRel(
	majStr string,
	minStr string,
	patchStr string,
) (int, int, int, PreReleaseType, error) {
	patchParts := strings.Split(patchStr, "-")

	var relVer PreReleaseType

	switch len(patchParts) {
	case 1:
		relVer = PreReleaseTypeNone
	case 2:
		patchStr = patchParts[0]
		switch patchParts[1] {
		case "alpha":
			relVer = PreReleaseTypeAlpha
		case "beta":
			relVer = PreReleaseTypeBeta
		}
	default:
		return 0, 0, 0, 0, fmt.Errorf("patch version %q: %w", patchStr, ErrMalformedVersion)
	}

	maj, min, patch, err := parseMajMinPatch(majStr, minStr, patchStr)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("error parsing version: %w", err)
	}

	return maj, min, patch, relVer, nil
}

func Parse(str string) (Version, error) {
	parts := strings.Split(str, ".")

	switch len(parts) {
	case 3:
		maj, min, patch, ver, err := parseRel(parts[0], parts[1], parts[2])
		if err != nil {
			return Version{}, fmt.Errorf("version %q: %w", str, err)
		}

		return Version{
			Version: MajMinPatch{
				Major: maj,
				Minor: min,
				Patch: patch,
			},
			PreRelease: PreRelease{
				Type: ver,
			},
		}, nil
	case 6:
		maj, min, patch, ver, err := parseRel(parts[0], parts[1], parts[2])
		if err != nil {
			return Version{}, fmt.Errorf("version %q: %w", str, err)
		}

		preMaj, preMin, prePatch, err := parseMajMinPatch(parts[3], parts[4], parts[5])
		if err != nil {
			return Version{}, fmt.Errorf("version %q: %w", str, err)
		}

		return Version{
			Version: MajMinPatch{
				Major: maj,
				Minor: min,
				Patch: patch,
			},
			PreRelease: PreRelease{
				Type: ver,
				Version: MajMinPatch{
					Major: preMaj,
					Minor: preMin,
					Patch: prePatch,
				},
			},
		}, nil
	}

	return Version{}, fmt.Errorf("version %q: %w", str, ErrMalformedVersion)
}
