package re

import (
	"regexp"
)

type Matcher struct {
	*regexp.Regexp
}

func MustCompile(expr string) Matcher {
	return Matcher{regexp.MustCompile(expr)}
}

func (m Matcher) Match(input string) map[string]string {
	match := m.FindStringSubmatch(input)
	if match == nil {
		return nil
	}
	result := make(map[string]string)
	for i, name := range m.SubexpNames() {
		result[name] = match[i]
	}
	return result
}
