package internal

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNotFound = errors.New("rule not found")
var ErrServiceExist = errors.New("service exist")

type Matcher interface {
	Match(req *http.Request) (string, error)
	Add(service string, rule string) error
}

type defaultMatcher struct {
	rules map[string]string
}

func (m *defaultMatcher) Add(service string, rule string) error {
	if _, ok := m.rules[service]; ok {
		return ErrServiceExist
	}

	m.rules[service] = rule
	return nil
}

func (m *defaultMatcher) Match(req *http.Request) (string, error) {
	for srv, rule := range m.rules {
		if strings.Contains(req.Host, rule) { //return true if host is xxx.domain.com, rule is domain.com
			return srv, nil
		}
	}
	return "", ErrNotFound
}

func NewMatcher() Matcher {
	return &defaultMatcher{rules: make(map[string]string)}
}
