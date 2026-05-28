package validator_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/thumbrise/commitlint-scope/pkg/validator"
)

func TestDefaultOutsiderFinder_Find(t *testing.T) {
	tests := []struct {
		name          string
		patterns      map[string][]string
		scope         string
		files         []string
		wantOutsiders []string
	}{
		{
			name: "no explicit patterns - uses scope/** as default",
			patterns: map[string][]string{
				"api": {},
			},
			scope:         "api",
			files:         []string{"api/handler.go", "api/helper.go", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name:          "no explicit patterns, even no configured scopes - uses scope/** as default",
			scope:         "api",
			files:         []string{"api/handler.go", "api/helper.go", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name: "exact match with wildcard",
			patterns: map[string][]string{
				"api": {"api/*"},
			},
			scope:         "api",
			files:         []string{"api/handler.go", "api/helper.go", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name: "globstar matches nested files",
			patterns: map[string][]string{
				"api": {"api/**"},
			},
			scope:         "api",
			files:         []string{"api/handler.go", "api/sub/helper.go", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name: "wildcard matches dotfiles (like dot:true)",
			patterns: map[string][]string{
				"env": {"*"},
			},
			scope:         "env",
			files:         []string{".env", "config.go"},
			wantOutsiders: []string{}, // * matches both
		},
		{
			name: "explicit dot pattern still works",
			patterns: map[string][]string{
				"env": {".*", "*.go"},
			},
			scope:         "env",
			files:         []string{".env", "config.go", ".gitignore"},
			wantOutsiders: []string{},
		},
		{
			name: "globstar matches dot directories",
			patterns: map[string][]string{
				"api": {"api/**"},
			},
			scope:         "api",
			files:         []string{"api/.internal/secret.go", "api/outer/file.go"},
			wantOutsiders: []string{},
		},
		{
			name: "complex patterns with multiple stars",
			patterns: map[string][]string{
				"db": {"db/migrations/*.sql", "db/schema/**"},
			},
			scope:         "db",
			files:         []string{"db/migrations/001_init.sql", "db/schema/latest.json", "db/docs/readme.md"},
			wantOutsiders: []string{"db/docs/readme.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := validator.NewDefaultOutsiderFinder(tt.patterns)
			got := finder.Find(tt.scope, tt.files)

			// Convert got outsidersFiles to slice of file names
			var gotFiles []string
			for _, o := range got {
				gotFiles = append(gotFiles, o.File)
			}

			if !sameSlice(t, gotFiles, tt.wantOutsiders) {
				t.Errorf("outsidersFiles = %v, want %v", gotFiles, tt.wantOutsiders)
			}
		})
	}
}

func sameSlice[T cmp.Ordered](t *testing.T, a, b []T) bool {
	t.Helper()

	if len(a) != len(b) {
		return false
	}

	slices.Sort(a)
	slices.Sort(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
