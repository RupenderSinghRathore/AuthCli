package main

import "errors"

var (
	ErrNotLoggedIn           = errors.New("not logged in")
	ErrUserNotFound          = errors.New("no such user")
	ErrIncorrectCode         = errors.New("wrong totp code")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)
