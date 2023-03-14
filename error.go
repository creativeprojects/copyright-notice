package main

import "fmt"

// ErrorClass for error types
type ErrorClass int

// ErrorClass
const (
	ErrorGeneric ErrorClass = iota
	FileErrorInvalidName
	FileErrorTooBig
	FileErrorCannotOpen
	FileErrorReading
)

func (e ErrorClass) String() string {
	switch e {
	case FileErrorInvalidName:
		return "invalid file descriptor"
	case FileErrorTooBig:
		return "file is too big"
	case FileErrorCannotOpen:
		return "cannot open file"
	case FileErrorReading:
		return "error reading file"
	default:
		return "error"
	}
}

// Error allows for error wrapping
type Error struct {
	class ErrorClass
	wrap  error
}

// NewError creates a new wrapped error
func NewError(class ErrorClass, wrap error) *Error {
	return &Error{
		class: class,
		wrap:  wrap,
	}
}

// Class returns the class of the error
func (e *Error) Class() ErrorClass {
	return e.class
}

func (e *Error) Error() string {
	if e.wrap == nil {
		return e.Class().String()
	}
	return fmt.Sprintf("%s: %s", e.class.String(), e.wrap.Error())
}

func (e *Error) Unwrap() error {
	return e.wrap
}

// verify interface
var _ error = &Error{}
