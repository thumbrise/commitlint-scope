package cmd

import (
	"github.com/thumbrise/commitlint-scope/cmd/commands"
	"github.com/urfave/cli/v3"
)

// Root is the entry point command for commitlint-scope.
var Root = &cli.Command{
	Name:        "commitlint-scope",
	Description: `commitlint-scope - a linter that checks if declared commit scopes match the changed files`,
	Commands: []*cli.Command{
		commands.RunCMD,
	},
	Suggest: true,
}
