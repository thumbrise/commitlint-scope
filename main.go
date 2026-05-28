package main

import (
	"context"
	"fmt"
	"os"

	"github.com/thumbrise/commitlint-scope/cmd"
)

func main() {
	ctx := context.Background()

	if err := cmd.Root.Run(ctx, os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal error: %[1]v\n", err)

		os.Exit(1)
	}
}
