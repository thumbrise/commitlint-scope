package validator_test

import (
	"regexp"
	"testing"

	"github.com/thumbrise/commitlint-scope/v2/pkg/validator"
)

func TestDefaultScopeParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		regex   *regexp.Regexp
		message string
		want    string
	}{
		{name: "nil regex", regex: nil, message: "feat(api): add", want: ""},
		{name: "no match", regex: regexp.MustCompile(`^feat\((?P<scope>[^)]+)\)`), message: "fix(api): bug", want: ""},
		{name: "match with scope", regex: regexp.MustCompile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`), message: "feat(api): add endpoint", want: "api"},
		{name: "no scope in message", regex: regexp.MustCompile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`), message: "chore: update deps", want: ""},
		{name: "regex without named group scope", regex: regexp.MustCompile(`^[a-z]+(?:\(([^)]+)\))?!?:\s`), message: "feat(api): add", want: ""},
		{name: "scope with breaking change (!)", regex: regexp.MustCompile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`), message: "fix(auth)!: correct token", want: "auth"},
		{name: "multiple matches, returns first", regex: regexp.MustCompile(`(?P<scope>[a-z]+)`), message: "api handler", want: "api"},
		{name: "empty scope in parentheses", regex: regexp.MustCompile(`^feat\((?P<scope>[^)]*)\)`), message: "feat(): empty", want: ""},
		{name: "composed scope", regex: regexp.MustCompile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`), message: "style(services,frontend): Add linters and formatters, format whole frontend", want: "services,frontend"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validator.NewDefaultScopeParser(tt.regex)

			got := p.Parse(tt.message)
			if got != tt.want {
				t.Errorf("Parse(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}
