package validator_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/thumbrise/commitlint-scope/v3/pkg/validator"
)

func TestDefaultScopeParser_Parse(t *testing.T) {
	defaultRegex := regexp.MustCompile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`)

	tests := []struct {
		name      string
		regex     *regexp.Regexp
		separator string
		message   string
		want      []string
	}{
		{name: "nil regex", regex: nil, message: "feat(api): add", want: nil},
		{name: "no match", regex: regexp.MustCompile(`^feat\((?P<scope>[^)]+)\)`), message: "fix(api): bug", want: nil},
		{name: "match with scope", regex: defaultRegex, message: "feat(api): add endpoint", want: []string{"api"}},
		{name: "no scope in message", regex: defaultRegex, message: "chore: update deps", want: nil},
		{name: "regex without named group scope", regex: regexp.MustCompile(`^[a-z]+(?:\(([^)]+)\))?!?:\s`), message: "feat(api): add", want: nil},
		{name: "scope with breaking change (!)", regex: defaultRegex, message: "fix(auth)!: correct token", want: []string{"auth"}},
		{name: "empty scope in parentheses", regex: regexp.MustCompile(`^feat\((?P<scope>[^)]*)\)`), message: "feat(): empty", want: nil},
		{name: "composed scope with comma", regex: defaultRegex, message: "style(services, frontend): Add linters", want: []string{"services", "frontend"}},
		{name: "composed scope with pipe separator", regex: defaultRegex, separator: "|", message: "style(services | frontend): Add linters", want: []string{"services", "frontend"}},
		{name: "multi scope with trailing comma", regex: defaultRegex, message: "feat(api,): add", want: []string{"api"}},
		{name: "triple scope", regex: defaultRegex, message: "feat(api, db, cache): add", want: []string{"api", "db", "cache"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validator.NewDefaultScopeParser(tt.regex, tt.separator)

			got := p.Parse(tt.message)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%q) = %v, want %v", tt.message, got, tt.want)
			}
		})
	}
}
