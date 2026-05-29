package validator_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/thumbrise/commitlint-scope/v2/pkg/validator"
)

func gitCall(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.CommandContext(t.Context(), "git", args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s\n%s", args, err, out)
	}

	return string(out)
}

func initRepo(t *testing.T) (string, []string) {
	t.Helper()

	dir := t.TempDir()

	gitCall(t, dir, "init")
	gitCall(t, dir, "config", "user.email", "test@test.com")
	gitCall(t, dir, "config", "user.name", "Test")

	writeFile(t, dir, "a.txt", "hello")
	gitCall(t, dir, "add", "a.txt")
	gitCall(t, dir, "commit", "-m", "first commit")

	writeFile(t, dir, "b.txt", "world")
	gitCall(t, dir, "add", "b.txt")
	gitCall(t, dir, "commit", "-m", "second commit")

	out := gitCall(t, dir, "rev-list", "--reverse", "HEAD")

	shas := splitLines(out)
	if len(shas) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(shas))
	}

	return dir, shas
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func splitLines(data string) []string {
	var lines []string

	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}

func TestDefaultGit_SHA(t *testing.T) {
	dir, shas := initRepo(t)
	g := validator.NewDefaultGit(dir)

	from := shas[0]
	to := shas[1]

	got, err := g.SHA(context.Background(), from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 || got[0] != shas[1] {
		t.Errorf("expected [%s], got %v", shas[1], got)
	}

	got, err = g.SHA(context.Background(), shas[1], shas[1])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestDefaultGit_Message(t *testing.T) {
	dir, shas := initRepo(t)
	g := validator.NewDefaultGit(dir)

	msg, err := g.Message(context.Background(), shas[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg != "first commit" {
		t.Errorf("expected 'first commit', got '%s'", msg)
	}

	msg, err = g.Message(context.Background(), shas[1])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg != "second commit" {
		t.Errorf("expected 'second commit', got '%s'", msg)
	}
}

func TestDefaultGit_FilesChanged(t *testing.T) {
	dir, shas := initRepo(t)
	g := validator.NewDefaultGit(dir)

	files, err := g.FilesChanged(context.Background(), shas[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 1 || files[0] != "a.txt" {
		t.Errorf("expected [a.txt], got %v", files)
	}

	files, err = g.FilesChanged(context.Background(), shas[1])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 1 || files[0] != "b.txt" {
		t.Errorf("expected [b.txt], got %v", files)
	}
}

func TestDefaultGit_ContextCancellation(t *testing.T) {
	dir, _ := initRepo(t)
	g := validator.NewDefaultGit(dir)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := g.SHA(ctx, "HEAD", "HEAD~1")
	if err == nil {
		t.Error("expected error due to cancelled context")
	}
}
