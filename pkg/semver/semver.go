package semver

import "fmt"

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
