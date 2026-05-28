package main

import (
	"context"
	"fmt"
	"os"

	"github.com/thumbrise/commitlint-scope/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ctx := context.Background()

	root := cmd.NewRootCMD(cmd.RootCMDLDFlags{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
	if err := root.Run(ctx, os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal error: %[1]v\n", err)

		os.Exit(1)
	}
}
