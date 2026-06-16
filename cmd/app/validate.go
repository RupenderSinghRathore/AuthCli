package main

import (
	"errors"
	"regexp"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}

	if !usernameRegex.MatchString(username) {
		return errors.New(
			"username must be 3-32 characters and contain only letters, numbers, '_' or '-'",
		)
	}

	return nil
}

func validatePassword(password []byte) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if len(password) > 72 {
		return errors.New("password must be at most 72 characters")
	}

	return nil
}
