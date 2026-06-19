package main

import (
	"regexp"

	"RupenderSinghRathore/AuthCli/internal/validator"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

func validateUsernamePassword(username string, password []byte) error {
	v := validator.New()
	validateUsername(v, username)
	validatePassword(v, password)

	if !v.Valid() {
		return v
	}
	return nil
}

func validateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "required")
	if username != "" {
		v.Check(
			usernameRegex.MatchString(username),
			"username",
			"must be 3-32 characters and contain only letters, numbers, '_' or '-'",
		)
	}
}

func validatePassword(v *validator.Validator, password []byte) {
	v.Check(len(password) >= 4, "password", "must be at least 4 characters")
	v.Check(len(password) <= 72, "password", "must be at most 72 characters")
}
