package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var EmailRX = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

type Validator struct {
	Errors []ValidationError
}

func (e Validator) Error() string {
	parts := make([]string, 0, len(e.Errors)+1)
	parts = append(parts, "Validation errors:")
	for _, err := range e.Errors {
		parts = append(parts, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(parts, "\n")
}

type ValidationError struct {
	Field   string
	Message string
}

func New() *Validator {
	return &Validator{Errors: make([]ValidationError, 0)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	v.Errors = append(v.Errors, ValidationError{Field: key, Message: message})
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
