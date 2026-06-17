package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"RupenderSinghRathore/AuthCli/internal/database"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotLoggedIn  = errors.New("not logged in")
	ErrUserNotFound = errors.New("no such user")
)

type WrongPasswordErr struct {
	Msg string
}

func (e WrongPasswordErr) Error() string {
	return e.Msg
}

func newWrongPassErr(failedAttempts int64) WrongPasswordErr {
	return WrongPasswordErr{
		Msg: fmt.Sprintf(
			"%d failed attempts, remaining %d",
			failedAttempts,
			MaxFailedAttempts-failedAttempts,
		),
	}
}

type AccountLockedErr struct {
	Until time.Time
}

func (e AccountLockedErr) Error() string {
	return fmt.Sprintf(
		"account locked until %s",
		e.Until.Format("2006-01-02 15:04:05"),
	)
}

func (app *application) createUser(username string, password []byte) (*database.User, error) {
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

func (app *application) tryUserLogin(username string, password []byte) (*database.User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	user, err := app.queary.GetUserByUsername(context.Background(), username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrUserNotFound
	case user.LockedUntil != nil && time.Now().Before(*user.LockedUntil):
		return nil, AccountLockedErr{Until: *user.LockedUntil}
	case err != nil:
		return nil, err
	}

	if !VerifyPassword(user.PasswordHash, string(password)) {
		// record unsuccessfull login
		lockedUntil := time.Now().Add(LockedUntil)
		user, err := app.queary.RecordFailedLogin(
			context.Background(),
			database.RecordFailedLoginParams{
				UserID:      user.ID,
				LockedUntil: &lockedUntil,
				MaxAttempts: MaxFailedAttempts,
			},
		)
		if err != nil {
			return nil, err
		}
		return nil, newWrongPassErr(user.FailedAttempts)
	}

	if err := app.queary.RecordSuccessfulLogin(context.Background(), user.ID); err != nil {
		return nil, err
	}

	return &user, nil
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
