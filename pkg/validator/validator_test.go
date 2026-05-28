package validator_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thumbrise/commitlint-scope/pkg/validator"
)

type TestableCommit struct {
	sha, message   string
	files          []string
	scope          string
	outsidersFiles []string
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
			scope:          commitScope,
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
		name           string
		commits        []TestableCommit
		wantViolations int
		wantErr        bool
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
				{sha: "bbb2222bbb", message: "feat(api): add endpoint", files: []string{"api/handler.go"}, scope: "api"},
			},
		},
		{
			name: "TestableCommit with outsider file",
			commits: []TestableCommit{
				{sha: "ccc3333ccc", message: "fix(auth): token", files: []string{"auth/service.go", "api/handler.go"}, scope: "auth", outsidersFiles: []string{"api/handler.go"}},
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
				{sha: "eee5555eee", message: "feat(ui): button", files: []string{}, scope: "ui"},
			},
		},
		{
			name: "git message error",
			commits: []TestableCommit{
				{sha: "fff6666fff", message: "", messageErr: errors.New("git command failed")},
			},
			wantErr: true,
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

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(violations) != tt.wantViolations {
				t.Errorf("got %d violations, want %d", len(violations), tt.wantViolations)
			}
		})
	}
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

		if c.scope != "" {
			parser.EXPECT().Parse(c.message).Return(c.scope)
		} else {
			parser.EXPECT().Parse(c.message).Return("")

			continue
		}

		git.EXPECT().FilesChanged(mock.Anything, c.sha).Return(c.files, nil)

		if len(c.files) == 0 {
			continue
		}

		outsiders := make([]validator.Outsider, len(c.outsidersFiles))
		for i, outsiderFile := range c.outsidersFiles {
			outsiders[i] = validator.Outsider{
				File:              outsiderFile,
				UnmatchedPatterns: nil,
			}
		}

		outsider.EXPECT().Find(c.scope, c.files).Return(outsiders)
	}
}
