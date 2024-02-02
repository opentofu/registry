package validate

import (
	"errors"
	"strings"
)

var ErrorEmptyList = errors.New("found empty list of versions")

type Errors []error

func (e Errors) Error() string {
	var buf strings.Builder
	for _, el := range e {
		_, _ = buf.WriteString(el.Error())
		_, _ = buf.WriteString("\n")
	}
	return buf.String()
}
