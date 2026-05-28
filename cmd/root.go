package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/fatih/color"
	"github.com/thumbrise/commitlint-scope/cmd/commands"
	"github.com/urfave/cli/v3"
)

// Root is the entry point command for commitlint-scope.
var Root = &cli.Command{
	Name:        "commitlint-scope",
	Description: `commitlint-scope - a linter that checks if declared commit scopes match the changed files`,
	Commands: []*cli.Command{
		commands.RunCMD,
		commands.InitCMD,
	},
	Suggest: true,
	ExitErrHandler: func(ctx context.Context, command *cli.Command, err error) {
		code := 1

		if coder, ok := errors.AsType[cli.ExitCoder](err); ok {
			code = coder.ExitCode()
		}

		_, _ = color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "\n%s\n", err)
		_, _ = color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "\nexit code %d\n", code)

		os.Exit(code)
	},
}
