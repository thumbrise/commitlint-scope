package validator

import (
	"regexp"
	"strings"
)

const scopeRegexGroup = "scope"

type DefaultScopeParser struct {
	scopeRegex *regexp.Regexp
	separator  string
}

func NewDefaultScopeParser(scopeRegex *regexp.Regexp, separator string) *DefaultScopeParser {
	if separator == "" {
		separator = ","
	}

	return &DefaultScopeParser{scopeRegex: scopeRegex, separator: separator}
}

func (p *DefaultScopeParser) Parse(message string) []string {
	regex := p.scopeRegex
	if regex == nil {
		return nil
	}

	matches := regex.FindStringSubmatch(message)
	if matches == nil {
		return nil
	}

	idx := regex.SubexpIndex(scopeRegexGroup)
	if idx < 0 || idx >= len(matches) {
		return nil
	}

	raw := matches[idx]
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, p.separator)

	scopes := make([]string, 0, len(parts))
	for _, part := range parts {
		s := strings.TrimSpace(part)
		if s != "" {
			scopes = append(scopes, s)
		}
	}

	if len(scopes) == 0 {
		return nil
	}

	return scopes
}
