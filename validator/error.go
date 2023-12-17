package validator

import (
	"fmt"
	"strings"
)

type validationError interface {
	Error() string
	Field() string
	Reason() string
}

type validationMultiError interface {
	Error() string
	AllErrors() []error
}

type ValidationError struct {
	field  string
	reason string
	cause  error
}

var _ validationError = ValidationError{}

func NewError(field, reason string) ValidationError {
	return ValidationError{field: field, reason: reason}
}

func NewErrorWithCause(field, reason string, cause error) ValidationError {
	return ValidationError{field: field, reason: reason, cause: cause}
}

func (e ValidationError) Field() string {
	return e.field
}

func (e ValidationError) Reason() string {
	return e.reason
}

func (e ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}
	return fmt.Sprintf("invalid %s: %s%s", e.field, e.reason, cause)
}

type ValidationMultiError []error

var _ validationMultiError = ValidationMultiError{}

func (m ValidationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

func (m ValidationMultiError) AllErrors() []error {
	return m
}

func NewMultiError(errs ...error) ValidationMultiError {
	return errs
}
