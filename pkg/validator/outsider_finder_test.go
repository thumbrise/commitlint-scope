package validator_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/thumbrise/commitlint-scope/v3/pkg/validator"
)

func TestDefaultOutsiderFinder_Find(t *testing.T) {
	tests := []struct {
		name          string
		patterns      []validator.OutsiderFinderPattern
		scope         string
		files         []string
		wantOutsiders []string
	}{
		{
			name: "no explicit patterns - uses scope/** as default",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"api"},
					Files:  []string{},
				},
			},
			scope:         "api",
			files:         []string{"api/handler.go", "api/helper.go", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name:          "no explicit patterns, even no configured scopes - uses scope/** as default",
			patterns:      nil,
			scope:         "api",
			files:         []string{"api/handler.go", "api/helper.go", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name: "exact match with wildcard",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"api"},
					Files:  []string{"api/*"},
				},
			},
			scope:         "api",
			files:         []string{"api/handler.go", "api/helper.go", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name: "globstar matches nested files",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"api"},
					Files:  []string{"api/**"},
				},
			},
			scope:         "api",
			files:         []string{"api/handler.go", "api/sub/helper.go", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name: "wildcard matches dotfiles (like dot:true)",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"env"},
					Files:  []string{"*"},
				},
			},
			scope:         "env",
			files:         []string{".env", "config.go"},
			wantOutsiders: []string{}, // * matches both
		},
		{
			name: "explicit dot pattern still works",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"env"},
					Files:  []string{".*", "*.go"},
				},
			},
			scope:         "env",
			files:         []string{".env", "config.go", ".gitignore"},
			wantOutsiders: []string{},
		},
		{
			name: "globstar matches dot directories",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"api"},
					Files:  []string{"api/**"},
				},
			},
			scope:         "api",
			files:         []string{"api/.internal/secret.go", "api/outer/file.go"},
			wantOutsiders: []string{},
		},
		{
			name: "complex patterns with multiple stars",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"db"},
					Files:  []string{"db/migrations/*.sql", "db/schema/**"},
				},
			},
			scope:         "db",
			files:         []string{"db/migrations/001_init.sql", "db/schema/latest.json", "db/docs/readme.md"},
			wantOutsiders: []string{"db/docs/readme.md"},
		},
		{
			name: "globstar at start matches root-level files (gobwas workaround)",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"deps"},
					Files:  []string{"**/go.mod"},
				},
			},
			scope:         "deps",
			files:         []string{"go.mod", "go.sum", "op/file.go", "subdir/go.mod"},
			wantOutsiders: []string{"go.sum", "op/file.go"},
		},
		{
			name: "globstar at start with wildcard matches root-level files",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"docs"},
					Files:  []string{"**/*.md"},
				},
			},
			scope:         "docs",
			files:         []string{"README.md", "CHANGELOG.md", "api/readme.md", "main.go"},
			wantOutsiders: []string{"main.go"},
		},
		{
			name: "globstar at start with dotted scope matches root-level file",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"v1.json"},
					Files:  []string{"**/rail.v1.json"},
				},
			},
			scope:         "v1.json",
			files:         []string{"rail.v1.json", "internal/rail.v1.json", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
		{
			name: "multiple scopes sharing same patterns",
			patterns: []validator.OutsiderFinderPattern{
				{
					Scopes: []string{"api", "v1.json"},
					Files:  []string{"**/rail.v1.json"},
				},
			},
			scope:         "v1.json",
			files:         []string{"internal/rail.v1.json", "external/rail.v1.json", "core/other.go"},
			wantOutsiders: []string{"core/other.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder, err := validator.NewDefaultOutsiderFinder(tt.patterns)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

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
