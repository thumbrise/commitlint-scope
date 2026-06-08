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
	ErrSkip            = errors.New("skip commit")
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
	Parse(message string) []string
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

		ofPatterns := make([]OutsiderFinderPattern, len(cfg.Patterns))
		for i, cfgPattern := range cfg.Patterns {
			ofPatterns[i] = OutsiderFinderPattern(cfgPattern)
		}

		outsiderFinder, err = NewDefaultOutsiderFinder(ofPatterns)
		if err != nil {
			return nil, err
		}
	}

	if scopeParser == nil {
		scopeParser = NewDefaultScopeParser(cfg.ScopeRegex, cfg.ScopeSeparator)
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
		violation, err := v.checkCommit(ctx, sha)
		if errors.Is(err, ErrSkip) {
			continue
		}

		if err != nil {
			return nil, err
		}

		violations = append(violations, *violation)
	}

	return violations, nil
}

func (v *Validator) checkCommit(ctx context.Context, sha string) (*Violation, error) {
	message, err := v.git.Message(ctx, sha)
	if err != nil {
		return nil, fmt.Errorf("%w sha=%s: %w", ErrGetMessage, sha, err)
	}

	if message == "" {
		v.logger.Debug("no message, skip", "sha", sha)

		return nil, fmt.Errorf("%w: empty message", ErrSkip)
	}

	scopes := v.scopeParser.Parse(message)
	if len(scopes) == 0 {
		v.logger.Debug("no scope, skip", "sha", sha, "message", message)

		return nil, fmt.Errorf("%w: no scope", ErrSkip)
	}

	files, err := v.git.FilesChanged(ctx, sha)
	if err != nil {
		return nil, fmt.Errorf("%w sha=%s, commit=%s: %w", ErrGetChangedFiles, sha, message, err)
	}

	if len(files) == 0 {
		v.logger.Debug("no files changed, skip", "sha", sha)

		return nil, fmt.Errorf("%w: no files", ErrSkip)
	}

	allOutsiders := findOutsiders(v.outsiderFinder, scopes, files)
	if len(allOutsiders) == 0 {
		return nil, fmt.Errorf("%w: no outsiders", ErrSkip)
	}

	truncatedSHA := sha
	if len(truncatedSHA) > v.shaLength {
		truncatedSHA = truncatedSHA[:v.shaLength]
	}

	return &Violation{
		SHA:       truncatedSHA,
		Header:    message,
		Outsiders: allOutsiders,
	}, nil
}

func findOutsiders(finder OutsiderFinder, scopes []string, files []string) []Outsider {
	if len(scopes) == 0 {
		return nil
	}

	outsidersByFile := make(map[string]Outsider)
	for _, o := range finder.Find(scopes[0], files) {
		outsidersByFile[o.File] = o
	}

	for _, scope := range scopes[1:] {
		validForScope := make(map[string]bool)
		for _, o := range finder.Find(scope, files) {
			validForScope[o.File] = true
		}

		for file := range outsidersByFile {
			if !validForScope[file] {
				delete(outsidersByFile, file)
			}
		}
	}

	result := make([]Outsider, 0, len(outsidersByFile))
	for _, o := range outsidersByFile {
		result = append(result, o)
	}

	return result
}
