package validator

import (
	"github.com/gobwas/glob"
)

// DefaultOutsiderFinder checks whether files match predefined glob patterns for a given scope.
type DefaultOutsiderFinder struct {
	humanFilePatterns map[string][]string
	globFilePatterns  map[string][]glob.Glob
}

// NewDefaultOutsiderFinder creates a new DefaultOutsiderFinder.
// Patterns is a map from scope name to a list of glob pattern strings.
// Patterns are compiled immediately; invalid patterns will panic.
func NewDefaultOutsiderFinder(patterns map[string][]string) *DefaultOutsiderFinder {
	globMap := make(map[string][]glob.Glob, len(patterns))
	for scope, pats := range patterns {
		globs := make([]glob.Glob, 0, len(pats))
		for _, p := range pats {
			globs = append(globs, glob.MustCompile(p, '/'))
		}

		globMap[scope] = globs
	}

	return &DefaultOutsiderFinder{
		humanFilePatterns: patterns,
		globFilePatterns:  globMap,
	}
}

// Find returns files that do not match any pattern for the given scope.
// If the scope has no patterns, zero config fallback applied - (scope + /**)
func (f *DefaultOutsiderFinder) Find(scope string, files []string) []Outsider {
	globFilePatterns := f.globFilePatterns[scope]
	humanFilePatterns := f.humanFilePatterns[scope]

	// zero config
	if len(globFilePatterns) == 0 {
		defaultPattern := scope + "/**"
		g := glob.MustCompile(defaultPattern, '/')
		globFilePatterns = []glob.Glob{g}
		humanFilePatterns = []string{defaultPattern}
	}

	var outsiders []Outsider

	for _, file := range files {
		if f.matchesAny(file, globFilePatterns) {
			continue
		}

		outsiders = append(outsiders, Outsider{
			File:              file,
			UnmatchedPatterns: humanFilePatterns,
		})
	}

	return outsiders
}

// matchesAny reports whether name matches at least one compiled glob.
func (f *DefaultOutsiderFinder) matchesAny(name string, globs []glob.Glob) bool {
	for _, g := range globs {
		if g.Match(name) {
			return true
		}
	}

	return false
}

// Outsider holds a file that failed the scope check and the patterns that were tested.
type Outsider struct {
	File              string   `json:"file,omitempty"`
	UnmatchedPatterns []string `json:"unmatchedPatterns,omitempty"`
}
