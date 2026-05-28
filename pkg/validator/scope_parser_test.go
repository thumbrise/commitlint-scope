package validator_test

import (
	"regexp"
	"testing"

	"github.com/thumbrise/commitlint-scope/pkg/validator"
)

func TestDefaultScopeParser_Parse(t *testing.T) {
	compile := regexp.MustCompile

	tests := []struct {
		name    string
		regex   *regexp.Regexp
		message string
		want    string
	}{
		{name: "nil regex", regex: nil, message: "feat(api): add", want: ""},
		{name: "no match", regex: compile(`^feat\((?P<scope>[^)]+)\)`), message: "fix(api): bug", want: ""},
		{name: "match with scope", regex: compile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`), message: "feat(api): add endpoint", want: "api"},
		{name: "no scope in message", regex: compile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`), message: "chore: update deps", want: ""},
		{name: "regex without named group scope", regex: compile(`^[a-z]+(?:\(([^)]+)\))?!?:\s`), message: "feat(api): add", want: ""},
		{name: "scope with breaking change (!)", regex: compile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`), message: "fix(auth)!: correct token", want: "auth"},
		{name: "multiple matches, returns first", regex: compile(`(?P<scope>[a-z]+)`), message: "api handler", want: "api"},
		{name: "empty scope in parentheses", regex: compile(`^feat\((?P<scope>[^)]*)\)`), message: "feat(): empty", want: ""},
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
