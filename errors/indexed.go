package errors

import "fmt"

type indexedError struct {
	cause error
	index int
}

// WithIndex attaches an array or slice index to the error.
// The data can be retreived by asserting the `Index` and `Cause` behaviors on the error.
func WithIndex(err error, i int) error {
	return indexedError{cause: err, index: i}
}

// Error implements the `error` interface.
func (e indexedError) Error() string {
	return fmt.Sprintf("[%d]: %v", e.index, e.cause)
}

// Index returns the contained array or slice index where the error occurred.
func (e indexedError) Index() int {
	return e.index
}

// Cause returns the contained error.
func (e indexedError) Cause() error {
	return e.cause
}
