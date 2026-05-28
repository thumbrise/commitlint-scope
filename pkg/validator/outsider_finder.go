package validator

import (
	"errors"
	"fmt"

	"github.com/gobwas/glob"
)

var ErrInvalidGlobPattern = errors.New("invalid glob pattern")

// DefaultOutsiderFinder checks whether files match predefined glob patterns for a given scope.
type DefaultOutsiderFinder struct {
	scopesToPatternsHuman map[string][]string
	scopesToPatternsGlob  map[string][]glob.Glob
}

// NewDefaultOutsiderFinder creates a new DefaultOutsiderFinder.
// scopesToPatterns is a map from scope name to a list of glob pattern strings.
func NewDefaultOutsiderFinder(scopesToPatterns map[string][]string) (*DefaultOutsiderFinder, error) {
	scopesToPatternsGlob := make(map[string][]glob.Glob, len(scopesToPatterns))
	for scope, patterns := range scopesToPatterns {
		globs := make([]glob.Glob, 0, len(patterns))
		for _, pattern := range patterns {
			g, err := glob.Compile(pattern, '/')
			if err != nil {
				return nil, fmt.Errorf("%w: %q: %w", ErrInvalidGlobPattern, pattern, err)
			}

			globs = append(globs, g)
		}

		scopesToPatternsGlob[scope] = globs
	}

	return &DefaultOutsiderFinder{
		scopesToPatternsHuman: scopesToPatterns,
		scopesToPatternsGlob:  scopesToPatternsGlob,
	}, nil
}

// Find returns files that do not match any pattern for the given scope.
// If the scope has no patterns, zero config fallback applied - (scope + /**)
func (f *DefaultOutsiderFinder) Find(scope string, files []string) []Outsider {
	globFilePatterns := f.scopesToPatternsGlob[scope]
	humanFilePatterns := f.scopesToPatternsHuman[scope]

	// zero config
	if len(globFilePatterns) == 0 {
		defaultPattern := scope + "/**"

		g, err := glob.Compile(defaultPattern, '/')
		if err != nil {
			errPanic := fmt.Sprintf("cannot compile default pattern %q: %s", defaultPattern, err.Error())
			panic(errPanic)
		}

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
