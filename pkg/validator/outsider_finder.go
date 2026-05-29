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

// OutsiderFinderPattern matches the new configuration structure.
type OutsiderFinderPattern struct {
	Scopes []string
	Files  []string
}

// NewDefaultOutsiderFinder creates a new DefaultOutsiderFinder.
// It accepts a slice of OutsiderFinderPattern and flattens it into efficient lookup maps.
func NewDefaultOutsiderFinder(patterns []OutsiderFinderPattern) (*DefaultOutsiderFinder, error) {
	scopesToPatternsHuman := make(map[string][]string)
	scopesToPatternsGlob := make(map[string][]glob.Glob)

	for _, item := range patterns {
		globs := make([]glob.Glob, 0, len(item.Files))
		for _, pattern := range item.Files {
			g, err := glob.Compile(pattern, '/')
			if err != nil {
				return nil, fmt.Errorf("%w: %q: %w", ErrInvalidGlobPattern, pattern, err)
			}

			globs = append(globs, g)
		}

		for _, scope := range item.Scopes {
			scopesToPatternsHuman[scope] = append(scopesToPatternsHuman[scope], item.Files...)
			scopesToPatternsGlob[scope] = append(scopesToPatternsGlob[scope], globs...)
		}
	}

	return &DefaultOutsiderFinder{
		scopesToPatternsHuman: scopesToPatternsHuman,
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
		defaultPattern := glob.QuoteMeta(scope) + "/**"

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
