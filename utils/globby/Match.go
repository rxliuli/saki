package globby

import (
	"github.com/gobwas/glob"
)

type MatchOptions struct {
	Ignore []string
}

func MiniMatch(patterns []glob.Glob, value string) bool {
	for _, pattern := range patterns {
		if pattern.Match(value) {
			return true
		}
	}
	return false
}

func Patterns2Globs(patterns []string) []glob.Glob {
	patternGlobs := make([]glob.Glob, 0)
	for _, pattern := range patterns {
		patternGlobs = append(patternGlobs, glob.MustCompile(pattern, '/'))
	}
	return patternGlobs
}
