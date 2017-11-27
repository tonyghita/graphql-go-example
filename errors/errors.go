// Package errors contains ...
package errors

import (
	"bytes"
	goerrors "errors"
	"fmt"
)

// Errors ...
type Errors []error

// Err ...
func (e Errors) Err() error {
	if len(e) == 0 {
		return nil
	}

	return e
}

// Error ...
func (e Errors) Error() string {
	var buf bytes.Buffer

	if n := len(e); n == 1 {
		buf.WriteString("1 error: ")
	} else {
		fmt.Fprintf(&buf, "%d errors: ", n)
	}

	for i, err := range e {
		if i != 0 {
			buf.WriteString("; ")
		}

		buf.WriteString(err.Error())
	}

	return buf.String()
}

// New ...
// This is convenience method so we don't have to fight with package imports.
func New(message string) error {
	return goerrors.New(message)
}
