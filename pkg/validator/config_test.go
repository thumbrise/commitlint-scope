package validator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thumbrise/commitlint-scope/v3/pkg/validator"
)

func TestLoadConfig(t *testing.T) {
	defaultRegexStr := `^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`

	tests := []struct {
		name         string
		yaml         string
		wantRegexNil bool
		wantRegex    string
		wantPatterns []validator.PatternItem
		wantErr      error
	}{
		{
			name:      "no config file",
			wantRegex: defaultRegexStr,
		},
		{
			name: "patterns with default regex",
			yaml: `patterns:
  - scopes: ["api"]
    files: ["api/*"]
  - scopes: ["core"]
    files: ["core/**"]
`,
			wantRegex: defaultRegexStr,
			wantPatterns: []validator.PatternItem{
				{Scopes: []string{"api"}, Files: []string{"api/*"}},
				{Scopes: []string{"core"}, Files: []string{"core/**"}},
			},
		},
		{
			name: "custom scopeRegex only",
			yaml: `scopeRegex: '^(feat|fix):'
`,
			wantRegex: `^(feat|fix):`,
		},
		{
			name: "both patterns and custom scopeRegex",
			yaml: `scopeRegex: '^(feat|fix):'
patterns:
  - scopes: ["api"]
    files: ["api/*"]
`,
			wantRegex: `^(feat|fix):`,
			wantPatterns: []validator.PatternItem{
				{Scopes: []string{"api"}, Files: []string{"api/*"}},
			},
		},
		{
			name: "patterns with dots inside scopes",
			yaml: `patterns:
  - scopes: ["rail.v1.json"]
    files: ["**/rail.v1.json"]
`,
			wantRegex: defaultRegexStr,
			wantPatterns: []validator.PatternItem{
				{Scopes: []string{"rail.v1.json"}, Files: []string{"**/rail.v1.json"}},
			},
		},
		{
			name:    "malformed config",
			yaml:    "[[invalid",
			wantErr: validator.ErrConfigRead,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			origDir, err := os.Getwd()
			require.NoError(t, err)

			t.Cleanup(func() {
				require.NoError(t, os.Chdir(origDir))
			})

			require.NoError(t, os.Chdir(dir))

			if tt.yaml != "" {
				err = os.WriteFile(filepath.Join(dir, ".commitlint-scope.yaml"), []byte(tt.yaml), 0o644)
				require.NoError(t, err)
			}

			cfg, err := validator.LoadConfig()
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}

			require.NoError(t, err)

			if tt.wantRegexNil {
				assert.Nil(t, cfg.ScopeRegex)
			} else {
				require.NotNil(t, cfg.ScopeRegex)
				assert.Equal(t, tt.wantRegex, cfg.ScopeRegex.String())
			}

			assert.Equal(t, tt.wantPatterns, cfg.Patterns)
		})
	}
}
