package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/thumbrise/commitlint-scope/pkg/validator"
	"github.com/urfave/cli/v3"
)

var (
	ErrFileExists = errors.New("file already exists")
	ErrOpen       = errors.New("cannot open file")
	ErrWrite      = errors.New("cannot write content")
)

const InitConfigData = `#$schema: https://github.com/thumbrise/commitlint-scope/blob/main/docs/schema/config.json

# Scope parsing customization. Not required, if you follow common conventional header. In example: 'type!(scope): subject'
#scopeRegex: ^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s
patterns:
  "auth": [ "services/auth/**" ]
  "migrations": [ "database/migrations/*.sql" ]
  "frontend": [ "**/assets/**", "**/frontend/**" ]
  "docs": [ "**/*.md" ]
`

const (
	InitConfigFileName = validator.ConfigName + ".yml"
	InitConfigFileMode = 0o600
)

var InitCMD = &cli.Command{
	Name:  "init",
	Usage: "Initialize a config file",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		f, err := os.OpenFile(InitConfigFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, InitConfigFileMode)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				return fmt.Errorf("%w: %s", ErrFileExists, InitConfigFileName)
			}

			return fmt.Errorf("%w: %w", ErrOpen, err)
		}

		defer func() {
			if closeErr := f.Close(); closeErr != nil && err == nil {
				_, _ = fmt.Fprintf(cmd.ErrWriter, "%s: %s", ErrWrite, closeErr)
			}
		}()

		_, err = f.WriteString(InitConfigData)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrWrite, err)
		}

		return nil
	},
}
