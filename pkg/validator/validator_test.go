package validator_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thumbrise/commitlint-scope/v3/pkg/validator"
)

type TestableCommit struct {
	sha, message   string
	files          []string
	scopes         []string
	outsidersFiles []string            // expected final outsiders (intersection)
	scopeOutsiders map[string][]string // per-scope outsiders returned by Find mock
	messageErr     error
}

func TestValidator_OneViolation(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)

	const (
		commitSHA     = "abc1234abcd"
		commitScope   = "api"
		commitMessage = "feat(" + commitScope + "): add endpoint"
	)

	filesChanged := []string{"core/other.go"}
	commits := []TestableCommit{
		{
			sha:            commitSHA,
			message:        commitMessage,
			files:          filesChanged,
			scopes:         []string{commitScope},
			outsidersFiles: filesChanged,
			messageErr:     nil,
		},
	}

	git := validator.NewMockGit(t)
	parser := validator.NewMockScopeParser(t)
	outsider := validator.NewMockOutsiderFinder(t)

	SetupExpectations(t, commits, git, parser, outsider)

	const shaLength = 7

	v, err := validator.NewValidator(
		validator.Config{},
		validator.Options{
			Logger:         logger,
			SHALength:      shaLength,
			Git:            git,
			OutsiderFinder: outsider,
			ScopeParser:    parser,
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	violations, err := v.Validate(context.Background(), "main", "feature-branch")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].SHA != commitSHA[:shaLength] {
		t.Errorf("expected SHA abc1234, got %s", violations[0].SHA)
	}

	if len(violations[0].Outsiders) != 1 {
		t.Errorf("expected 1 outsider, got %d", len(violations[0].Outsiders))
	}
}

func TestValidator_Scenarios(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)

	const shaLength = 7

	tests := []struct {
		name              string
		commits           []TestableCommit
		wantViolations    int
		wantOutsiderFiles []string
		wantErr           bool
	}{
		{
			name: "TestableCommit without scope skipped",
			commits: []TestableCommit{
				{sha: "aaa1111aaa", message: "chore: update deps", files: []string{"go.mod"}},
			},
		},
		{
			name: "TestableCommit with scope and no outsidersFiles",
			commits: []TestableCommit{
				{sha: "bbb2222bbb", message: "feat(api): add endpoint", files: []string{"api/handler.go"}, scopes: []string{"api"}},
			},
		},
		{
			name: "TestableCommit with outsider file",
			commits: []TestableCommit{
				{sha: "ccc3333ccc", message: "fix(auth): token", files: []string{"auth/service.go", "api/handler.go"}, scopes: []string{"auth"}, outsidersFiles: []string{"api/handler.go"}},
			},
			wantViolations: 1,
		},
		{
			name: "empty message skipped",
			commits: []TestableCommit{
				{sha: "ddd4444ddd", message: "", files: []string{"file.go"}},
			},
		},
		{
			name: "no files changed skipped",
			commits: []TestableCommit{
				{sha: "eee5555eee", message: "feat(ui): button", files: []string{}, scopes: []string{"ui"}},
			},
		},
		{
			name: "git message error",
			commits: []TestableCommit{
				{sha: "fff6666fff", message: "", messageErr: errors.New("git command failed")},
			},
			wantErr: true,
		},
		{
			name: "multi scope with outsiders intersection",
			commits: []TestableCommit{
				{
					sha:     "ggg7777ggg",
					message: "feat(api, db): add endpoint",
					files:   []string{"api/handler.go", "db/schema.sql", "core/other.go", "README.md"},
					scopes:  []string{"api", "db"},
					scopeOutsiders: map[string][]string{
						"api": {"db/schema.sql", "core/other.go", "README.md"},
						"db":  {"api/handler.go", "core/other.go", "README.md"},
					},
					outsidersFiles: []string{"core/other.go", "README.md"},
				},
			},
			wantViolations:    1,
			wantOutsiderFiles: []string{"core/other.go", "README.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := validator.NewMockGit(t)
			parser := validator.NewMockScopeParser(t)
			outsider := validator.NewMockOutsiderFinder(t)

			SetupExpectations(t, tt.commits, git, parser, outsider)

			v, err := validator.NewValidator(
				validator.Config{},
				validator.Options{
					Logger:         logger,
					SHALength:      shaLength,
					Git:            git,
					OutsiderFinder: outsider,
					ScopeParser:    parser,
				})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			violations, err := v.Validate(context.Background(), "main", "feature-branch")

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(violations) != tt.wantViolations {
				t.Errorf("got %d violations, want %d", len(violations), tt.wantViolations)

				return
			}

			if len(tt.wantOutsiderFiles) > 0 {
				gotFiles := make([]string, 0, len(tt.wantOutsiderFiles))
				for _, o := range violations[0].Outsiders {
					gotFiles = append(gotFiles, o.File)
				}

				if !sameSlice(t, gotFiles, tt.wantOutsiderFiles) {
					t.Errorf("outsiders = %v, want %v", gotFiles, tt.wantOutsiderFiles)
				}
			}
		})
	}
}

func outsidersFromFiles(files []string) []validator.Outsider {
	out := make([]validator.Outsider, len(files))
	for i, f := range files {
		out[i] = validator.Outsider{File: f, UnmatchedPatterns: nil}
	}

	return out
}

func SetupExpectations(t *testing.T, commits []TestableCommit, git *validator.MockGit, parser *validator.MockScopeParser, outsider *validator.MockOutsiderFinder) {
	t.Helper()

	shas := make([]string, len(commits))
	for i, c := range commits {
		shas[i] = c.sha
	}

	git.EXPECT().SHA(mock.Anything, mock.Anything, mock.Anything).Return(shas, nil)

	for _, c := range commits {
		if c.messageErr != nil {
			git.EXPECT().Message(mock.Anything, c.sha).Return("", c.messageErr)

			continue
		}

		git.EXPECT().Message(mock.Anything, c.sha).Return(c.message, nil)

		if c.message == "" {
			continue
		}

		if len(c.scopes) > 0 {
			parser.EXPECT().Parse(c.message).Return(c.scopes)
		} else {
			parser.EXPECT().Parse(c.message).Return(nil)

			continue
		}

		git.EXPECT().FilesChanged(mock.Anything, c.sha).Return(c.files, nil)

		if len(c.files) == 0 {
			continue
		}

		for _, scope := range c.scopes {
			scopeFiles, ok := c.scopeOutsiders[scope]
			if !ok {
				scopeFiles = c.outsidersFiles
			}

			outsider.EXPECT().Find(scope, c.files).Return(outsidersFromFiles(scopeFiles))
		}
	}
}
