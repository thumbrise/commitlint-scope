package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/thumbrise/commitlint-scope/v3/cmd/errs"
	"github.com/thumbrise/commitlint-scope/v3/pkg/validator"
	"github.com/urfave/cli/v3"
)

var (
	flagRunFrom    string
	flagRunTo      string
	flagRunVerbose bool
	flagRunNoColor bool
	flagRunJSON    bool
)

var RunCMD = &cli.Command{
	Name:      "run",
	Usage:     "Lint commit scopes against changed files",
	UsageText: "commitlint-scope run --from <sha> --to <sha>",
	Description: `Validate that scopes declared in commit messages correspond to actually changed files

The command inspects a range of commits (from exclusive, to inclusive) and reports any scope that does not match the files modified in that commit.

Examples:
  commitlint-scope run --from main --to feature-branch
  commitlint-scope run --from HEAD~5 --to HEAD
`,
	Suggest: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "from",
			Aliases:     []string{"f"},
			Usage:       "Start of the commit range (exclusive)",
			Required:    true,
			Destination: &flagRunFrom,
		},
		&cli.StringFlag{
			Name:        "to",
			Aliases:     []string{"t"},
			Usage:       "End of the commit range (inclusive)",
			Required:    true,
			Destination: &flagRunTo,
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Usage:       "Verbose output",
			Required:    false,
			Destination: &flagRunVerbose,
		},
		&cli.BoolFlag{
			Name:        "no-color",
			Usage:       "Disable color output",
			Required:    false,
			Destination: &flagRunNoColor,
		},
		&cli.BoolFlag{
			Name:        "json",
			Usage:       "Show output in JSON format",
			Required:    false,
			Destination: &flagRunJSON,
		},
	},

	Action: func(ctx context.Context, cmd *cli.Command) error {
		color.NoColor = flagRunNoColor

		configureLogger(os.Stderr, flagRunVerbose)

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

		violations, err := vld.Validate(ctx, flagRunFrom, flagRunTo)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		if len(violations) == 0 {
			return nil
		}

		if flagRunJSON {
			err := jsonOutput(cmd.Writer, violations)
			if err != nil {
				return err
			}
		} else {
			textOutput(cmd.Writer, violations)
		}

		return errs.NewViolationsFound(len(violations))
	},
}

func textOutput(w io.Writer, violations []validator.Violation) {
	shaColor := color.New(color.FgYellow, color.Bold)
	for _, v := range violations {
		_, _ = fmt.Fprintf(w, "%s %s\n", shaColor.Sprintf("%s", v.SHA), v.Header)
		for _, o := range v.Outsiders {
			_, _ = fmt.Fprintf(w, "  - %s\n", o.File)
		}
	}
}

func jsonOutput(writer io.Writer, violations []validator.Violation) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(violations); err != nil {
		return fmt.Errorf("failed to output violations: %w", err)
	}

	return nil
}

func configureLogger(writer io.Writer, verbose bool) {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}

	if verbose {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}

	handler := slog.NewTextHandler(writer, opts)

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
