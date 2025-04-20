package filter

import (
	"fmt"
	"regexp"
)

type SecretFilter struct {
	patterns []*regexp.Regexp
	action   string
}

func NewSecretFilter(patterns []string, action string) (*SecretFilter, error) {
	var regexps []*regexp.Regexp
	for _, pat := range patterns {
		r, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern '%s': %w", pat, err)
		}
		regexps = append(regexps, r)
	}
	return &SecretFilter{patterns: regexps, action: action}, nil
}

func (sf *SecretFilter) FilterLine(line string) (string, bool) {
	for _, r := range sf.patterns {
		if r.MatchString(line) {
			if sf.action == "redact" {
				return "[REDACTED]", true
			} else if sf.action == "block" {
				return "", true
			}
		}
	}
	return line, false
}
