package validator

import "regexp"

const scopeRegexGroup = "scope"

type DefaultScopeParser struct {
	scopeRegex *regexp.Regexp
}

func NewDefaultScopeParser(scopeRegex *regexp.Regexp) *DefaultScopeParser {
	return &DefaultScopeParser{scopeRegex: scopeRegex}
}

func (p *DefaultScopeParser) Parse(message string) string {
	regex := p.scopeRegex
	if regex == nil {
		return ""
	}

	matches := regex.FindStringSubmatch(message)
	if matches == nil {
		return ""
	}

	idx := regex.SubexpIndex(scopeRegexGroup)
	if idx < 0 || idx >= len(matches) {
		return ""
	}

	return matches[idx]
}
