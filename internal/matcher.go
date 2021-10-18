package internal

import (
	"errors"
	"net/http"
	"strings"
)

//ErrNotFound service's rule not found
var ErrNotFound = errors.New("rule not found")

//ErrServiceExist service register already
var ErrServiceExist = errors.New("service exist")

//Matcher service matcher interface
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

//NewMatcher create matcher
func NewMatcher() Matcher {
	return &defaultMatcher{rules: make(map[string]string)}
}
