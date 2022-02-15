package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	MatcherSuffix = "SUFFIX"
	MatcherExact  = "EXACT"
	MatcherRegex  = "REGEX"
)

type DomainMatcher interface {
	IsMatch(query string) bool
}

func ParseMatcher(rule string) (DomainMatcher, error) {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return nil, errors.New("Empty pattern")
	}

	// deconstruct pattern
	parts := strings.Split(rule, "|")
	pattern := strings.TrimSpace(parts[0])
	typ := MatcherSuffix
	if len(parts) > 1 {
		typ = strings.ToUpper(strings.TrimSpace(parts[1]))
	}

	switch typ {
	case MatcherSuffix:
		return newSuffixMatcher(pattern), nil
	case MatcherExact:
		return newExactMatcher(pattern), nil
	case MatcherRegex:
		return newRegexMatcher(pattern)
	default:
		return nil, errors.New("Unknown pattern type")
	}
}

type suffixMatcher struct {
	fqdn string
}

func newSuffixMatcher(pattern string) *suffixMatcher {
	return &suffixMatcher{
		fqdn: fmt.Sprintf("%s.", pattern),
	}
}

func (m *suffixMatcher) IsMatch(query string) bool {
	return strings.HasSuffix(query, m.fqdn)
}

type exactMatcher struct {
	fqdn string
}

func newExactMatcher(pattern string) *exactMatcher {
	return &exactMatcher{
		fqdn: fmt.Sprintf("%s.", pattern),
	}
}

func (m *exactMatcher) IsMatch(query string) bool {
	return query == m.fqdn
}

type regexMatcher struct {
	pattern *regexp.Regexp
}

func newRegexMatcher(pattern string) (*regexMatcher, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &regexMatcher{
		pattern: r,
	}, nil
}

func (m *regexMatcher) IsMatch(query string) bool {
	return m.pattern.MatchString(query)
}
