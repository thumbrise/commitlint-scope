package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/thumbrise/commitlint-scope/pkg/validator"
	"github.com/urfave/cli/v3"
)

var (
	from string
	to   string
)

// Root is the entry point command for commitlint-scope.
var Root = &cli.Command{
	Name:      "commitlint-scope",
	Usage:     "Lint commit scopes against changed files.",
	UsageText: "commitlint-scope --from <sha> --to <sha>",
	Description: `Validate that scopes declared in commit messages correspond to actually changed files.

The command inspects a range of commits (from exclusive, to inclusive) and reports any scope that does not match the files modified in that commit.

Examples:
  commitlint-scope --from main --to feature-branch
  commitlint-scope --from HEAD~5 --to HEAD
`,

	Suggest: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "from",
			Aliases:     []string{"f"},
			Usage:       "Start of the commit range (exclusive)",
			Required:    true,
			Destination: &from,
		},
		&cli.StringFlag{
			Name:        "to",
			Aliases:     []string{"t"},
			Usage:       "End of the commit range (inclusive)",
			Required:    true,
			Destination: &to,
		},
	},

	Action: func(ctx context.Context, cmd *cli.Command) error {
		cfg, err := validator.LoadConfig()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		vld, err := validator.NewValidator(cfg, validator.Options{
			Logger:         slog.Default(),
			SHALength:      7,
			Git:            nil,
			OutsiderFinder: nil,
			ScopeParser:    nil,
		})
		if err != nil {
			return fmt.Errorf("creating validator: %w", err)
		}

		violations, err := vld.Validate(ctx, from, to)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		if len(violations) == 0 {
			return nil
		}

		encoder := json.NewEncoder(cmd.Writer)
		encoder.SetIndent("", "  ")

		for _, v := range violations {
			if err := encoder.Encode(v); err != nil {
				return fmt.Errorf("failed to output violation: %w", err)
			}
		}

		return nil
	},
}
