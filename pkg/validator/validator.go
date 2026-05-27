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
)

type Validator struct {
	logger         *slog.Logger
	git            Git
	outsiderFinder OutsiderFinder
	scopeParser    ScopeParser
	shaLength      int
}

func NewValidator(
	logger *slog.Logger,
	git Git,
	outsiderFinder OutsiderFinder,
	scopeParser ScopeParser,
	shaLength int,
) *Validator {
	if logger == nil {
		logger = slog.Default()
	}

	if git == nil {
		panic("git must not be nil")
	}

	if outsiderFinder == nil {
		panic("outsiderFinder must not be nil")
	}

	if scopeParser == nil {
		panic("scopeParser must not be nil")
	}

	if shaLength < 1 {
		panic("shaLength must be greater than or equal to 1")
	}

	return &Validator{
		logger:         logger,
		git:            git,
		outsiderFinder: outsiderFinder,
		scopeParser:    scopeParser,
		shaLength:      shaLength,
	}
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

		scope, ok := v.scopeParser.Parse(message)
		if !ok {
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
			violations = append(violations, Violation{
				SHA:       sha[:v.shaLength],
				Header:    message,
				Outsiders: outsiders,
			})
		}
	}

	return violations, nil
}
