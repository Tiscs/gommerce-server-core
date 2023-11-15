package validator

type validationError interface {
	Error() string
	Field() string
	Reason() string
}

type validationMultiError interface {
	Error() string
	AllErrors() []error
}
