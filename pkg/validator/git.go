package validator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type DefaultGit struct {
	dir string
}

func NewDefaultGit(baseDir string) *DefaultGit {
	return &DefaultGit{dir: baseDir}
}

func (d *DefaultGit) SHA(ctx context.Context, from, to string) ([]string, error) {
	rangeArg := fmt.Sprintf("%s..%s", from, to)
	cmd := d.gitCommand(ctx, "rev-list", rangeArg)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git rev-list %s: %w", rangeArg, err)
	}

	return splitLines(out), nil
}

func (d *DefaultGit) Message(ctx context.Context, sha string) (string, error) {
	cmd := d.gitCommand(ctx, "log", "--format=%s", "-1", sha)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git log %s: %w", sha, err)
	}

	return strings.TrimSpace(string(out)), nil
}

func (d *DefaultGit) FilesChanged(ctx context.Context, sha string) ([]string, error) {
	cmd := d.gitCommand(ctx, "diff-tree", "--no-commit-id", "-r", "--root", "--name-only", sha)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff-tree %s: %w", sha, err)
	}

	return splitLines(out), nil
}

func (d *DefaultGit) gitCommand(ctx context.Context, args ...string) *exec.Cmd {
	//nolint:gosec // no other way
	cmd := exec.CommandContext(ctx, "git", args...)
	if d.dir != "" {
		cmd.Dir = d.dir
	}

	return cmd
}

func splitLines(data []byte) []string {
	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return nil
	}

	lines := strings.Split(raw, "\n")

	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}

	return result
}
