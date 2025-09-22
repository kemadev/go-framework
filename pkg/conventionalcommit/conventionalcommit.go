package conventionalcommit

import (
	"errors"
	"fmt"
	"regexp"
)

var ErrMalformedCommit = errors.New("commit message is malformed")

type BumpType int

const (
	BumpTypeMajor BumpType = iota
	BumpTypeMinor
	BumpTypePatch
)

type Commit struct {
	Type         string
	Scope        string
	BreakingBang string
}

var MinorTags []string = []string{
	"feat",
}

var PatchTags []string = []string{
	"fix",
	"perf",
	"revert",
}

var NoBumpTags []string = []string{
	"build",
	"chore",
	"ci",
	"docs",
	"style",
	"test",
	"refactor",
}

const (
	CaptureGroupKeyType         string = "committype"
	CaptureGroupKeyScope        string = "commitscope"
	CaptureGroupKeyBreakingBang string = "breakingbang"
)

const CommitRegex string = `^(?P<` + CaptureGroupKeyType + `>[a-zA-Z]+)(?P<` + CaptureGroupKeyScope + `>\([[:alnum:]._-]+\))?(?P<` + CaptureGroupKeyBreakingBang + `>!)?: ([[:alnum:]])+([[:space:][:print:]])*`

func Parse(msg string) (Commit, error) {
	reg, err := regexp.Compile(CommitRegex)
	if err != nil {
		return Commit{}, fmt.Errorf("error compiling regex %q: %w", CommitRegex, err)
	}

	matches := reg.FindStringSubmatch(msg)
	if matches == nil {
		return Commit{}, fmt.Errorf("version %q: %w", msg, ErrMalformedCommit)
	}

	names := reg.SubexpNames()
	result := make(map[string]string)

	for i, name := range names {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	return Commit{
		Type:         result[CaptureGroupKeyType],
		Scope:        result[CaptureGroupKeyScope],
		BreakingBang: result[CaptureGroupKeyBreakingBang],
	}, nil
}
