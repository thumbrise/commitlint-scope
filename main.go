package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/thumbrise/commitlint-scope/cmd"
)

func main() {
	ctx := context.Background()

	configureLogger()

	if err := cmd.Root.Run(ctx, os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %[1]v\n", err)

		os.Exit(1)
	}
}

func configureLogger() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     nil,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
