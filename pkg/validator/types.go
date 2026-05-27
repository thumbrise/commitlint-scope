package validator

import "context"

type Violation struct {
	SHA       string   `json:"sha"`
	Header    string   `json:"header"`
	Outsiders []string `json:"outsiders"`
}
type Git interface {
	SHA(ctx context.Context, from, to string) ([]string, error)
	Message(ctx context.Context, sha string) (string, error)
	FilesChanged(ctx context.Context, sha string) ([]string, error)
}
type ScopeParser interface {
	Parse(message string) (string, bool) // scope, ok
}
type OutsiderFinder interface {
	Find(scope string, files []string) []string
}
