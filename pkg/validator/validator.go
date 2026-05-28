package validator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrGetMessage      = errors.New("get commit message")
	ErrGetChangedFiles = errors.New("get changed files")
	ErrShaLength       = errors.New("sha length must be greater than 0")
)

type Violation struct {
	SHA       string     `json:"sha"`
	Header    string     `json:"header"`
	Outsiders []Outsider `json:"outsiders"`
}
type Git interface {
	SHA(ctx context.Context, from, to string) ([]string, error)
	Message(ctx context.Context, sha string) (string, error)
	FilesChanged(ctx context.Context, sha string) ([]string, error)
}
type ScopeParser interface {
	Parse(message string) string
}
type OutsiderFinder interface {
	Find(scope string, files []string) []Outsider
}

type Options struct {
	Logger         *slog.Logger
	SHALength      int
	Git            Git
	OutsiderFinder OutsiderFinder
	ScopeParser    ScopeParser
}
type Validator struct {
	logger         *slog.Logger
	git            Git
	outsiderFinder OutsiderFinder
	scopeParser    ScopeParser
	shaLength      int
}

func NewValidator(cfg Config, options Options) (*Validator, error) {
	logger := options.Logger
	shaLength := options.SHALength
	scopeParser := options.ScopeParser
	outsiderFinder := options.OutsiderFinder
	git := options.Git

	if logger == nil {
		logger = slog.Default()
	}

	if git == nil {
		git = NewDefaultGit("")
	}

	if outsiderFinder == nil {
		var err error

		outsiderFinder, err = NewDefaultOutsiderFinder(cfg.Patterns)
		if err != nil {
			return nil, err
		}
	}

	if scopeParser == nil {
		scopeParser = NewDefaultScopeParser(cfg.ScopeRegex)
	}

	if shaLength == 0 {
		shaLength = 7
	}

	if shaLength < 0 {
		return nil, fmt.Errorf("%w, got %d", ErrShaLength, shaLength)
	}

	return &Validator{
		logger:         logger,
		git:            git,
		outsiderFinder: outsiderFinder,
		scopeParser:    scopeParser,
		shaLength:      shaLength,
	}, nil
}

func (v *Validator) Validate(ctx context.Context, from, to string) ([]Violation, error) {
	shas, err := v.git.SHA(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("git sha: %w", err)
	}

	var violations []Violation

	for _, sha := range shas {
		message, err := v.git.Message(ctx, sha)
		if err != nil {
			return nil, fmt.Errorf("%w sha=%s: %w", ErrGetMessage, sha, err)
		}

		if message == "" {
			v.logger.Info("no message, skip", "sha", sha)

			continue
		}

		scope := v.scopeParser.Parse(message)
		if scope == "" {
			v.logger.Info("no scope, skip", "sha", sha, "message", message)

			continue
		}

		files, err := v.git.FilesChanged(ctx, sha)
		if err != nil {
			return nil, fmt.Errorf("%w sha=%s, commit=%s: %w", ErrGetChangedFiles, sha, message, err)
		}

		if len(files) == 0 {
			v.logger.Info("no files changed, skip", "sha", sha)

			continue
		}

		outsiders := v.outsiderFinder.Find(scope, files)
		if len(outsiders) > 0 {
			truncatedSHA := sha
			if len(truncatedSHA) > v.shaLength {
				truncatedSHA = truncatedSHA[:v.shaLength]
			}

			violations = append(violations, Violation{
				SHA:       truncatedSHA,
				Header:    message,
				Outsiders: outsiders,
			})
		}
	}

	return violations, nil
}
