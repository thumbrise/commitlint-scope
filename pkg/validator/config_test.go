package validator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thumbrise/commitlint-scope/v2/pkg/validator"
)

func TestLoadConfig(t *testing.T) {
	defaultRegexStr := `^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`

	tests := []struct {
		name         string
		yaml         string
		wantRegexNil bool
		wantRegex    string
		wantPatterns map[string][]string
		wantErr      error
	}{
		{
			name:      "no config file",
			wantRegex: defaultRegexStr,
		},
		{
			name: "patterns with default regex",
			yaml: `patterns:
  api:
    - api/*
  core:
    - core/**
`,
			wantRegex: defaultRegexStr,
			wantPatterns: map[string][]string{
				"api":  {"api/*"},
				"core": {"core/**"},
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
  api:
    - api/*
`,
			wantRegex: `^(feat|fix):`,
			wantPatterns: map[string][]string{
				"api": {"api/*"},
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
