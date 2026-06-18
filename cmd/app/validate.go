package main

import (
	"regexp"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

type ValidatationErr struct {
	Msg string
}

func (e ValidatationErr) Error() string {
	return e.Msg
}

func newValidationErr(msg string) ValidatationErr {
	return ValidatationErr{Msg: msg}
}

func validateUsernamePassowrd(username string, password []byte) error {
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
		return newValidationErr("Err: password must be at least 8 characters")
	}

	if len(password) > 72 {
		return newValidationErr("Err: password must be at most 72 characters")
	}

	return nil
}

func (app *application) validateUsernamePassowrd(username string, password []byte) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	return nil
}
