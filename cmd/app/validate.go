package main

import (
	"regexp"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

type ValidationErr struct {
	Msg string
}

func (e ValidationErr) Error() string {
	return e.Msg
}

func newValidationErr(msg string) ValidationErr {
	return ValidationErr{Msg: msg}
}

func validateUsernamePassword(username string, password []byte) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return newValidationErr("Err: username is required")
	}

	if !usernameRegex.MatchString(username) {
		return newValidationErr(
			"Err: username must be 3-32 characters and contain only letters, numbers, '_' or '-'",
		)
	}

	return nil
}

func validatePassword(password []byte) error {
	if len(password) < 4 {
		return newValidationErr("Err: password must be at least 4 characters")
	}

	if len(password) > 72 {
		return newValidationErr("Err: password must be at most 72 characters")
	}

	return nil
}
