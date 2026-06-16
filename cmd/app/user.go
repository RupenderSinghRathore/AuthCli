package main

import (
	"context"

	"RupenderSinghRathore/AuthCli/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) register(username string, password []byte) (*database.User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	PasswordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err := app.queary.CreateUser(
		context.Background(),
		database.CreateUserParams{Username: username, PasswordHash: PasswordHash},
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (app *application) login(username string, password []byte) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	return nil
}

func (app *application) logout() {
}

func HashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(hash), err
}

func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}
