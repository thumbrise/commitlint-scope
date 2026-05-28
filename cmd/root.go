package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/thumbrise/commitlint-scope/cmd/commands"
	"github.com/urfave/cli/v3"
)

type RootCMDLDFlags struct {
	Version string
	Commit  string
	Date    string
}

type RootCMD struct {
	cmd *cli.Command
}

func NewRootCMD(ldflags RootCMDLDFlags) *RootCMD {
	return &RootCMD{
		cmd: &cli.Command{
			Name:    "commitlint-scope",
			Usage:   "linter that checks if declared commit scopes match the changed files",
			Version: fmt.Sprintf("v%s %s %s", ldflags.Version, ldflags.Commit, ldflags.Date),
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
		},
	}
}

func (c *RootCMD) Run(ctx context.Context, args []string) error {
	return c.cmd.Run(ctx, args)
}
