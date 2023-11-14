package validator

type validationError interface {
	Field() string
	Reason() string
	Cause() error
	Key() bool
	ErrorName() string
	Error() string
}

type validationMultiError interface {
	AllErrors() []error
	Error() string
}
