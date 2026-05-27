package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var (
	from string
	to   string
)

// Violation represents a single scope violation found in a commit.
type Violation struct {
	SHA       string   `json:"sha"`
	Header    string   `json:"header"`
	Outsiders []string `json:"outsiders"`
}

// ScopeValidator defines the contract for validating commit scopes.
// Implement this interface with your business logic.
type ScopeValidator interface {
	Validate(ctx context.Context, from, to string) ([]Violation, error)
}

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
		//nolint:godox // Need
		// TODO: Replace the following line with your real validator initialization.
		//  Example:
		//    validator := myvalidator.New(git.New(), parser.New())
		var validator ScopeValidator

		//nolint:godox // Need
		// TODO: Instantiate implementation of ScopeValidator here.
		//   validator = yourpkg.NewValidator(...)
		if validator == nil {
			fmt.Println("Validator not yet implemented. Please initialize the ScopeValidator in cmd/root.go")

			return nil
		}

		violations, err := validator.Validate(ctx, from, to)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		if len(violations) == 0 {
			return nil
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")

		for _, v := range violations {
			if err := encoder.Encode(v); err != nil {
				return fmt.Errorf("failed to output violation: %w", err)
			}
		}

		return nil
	},
}
