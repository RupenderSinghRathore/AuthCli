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
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err := app.queries.CreateUser(
		context.Background(),
		database.CreateUserParams{Username: username, PasswordHash: passwordHash},
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (app *application) getUserForLogin(username string) (*database.User, error) {
	user, err := app.queries.GetUserByUsername(context.Background(), username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrUserNotFound
	case user.LockedUntil != nil && time.Now().Before(*user.LockedUntil):
		return nil, AccountLockedErr{Until: *user.LockedUntil}
	case err != nil:
		return nil, err
	}

	return &user, nil
}

func (app *application) recordSuccessfulLogin(user *database.User) error {
	return app.queries.RecordSuccessfulLogin(context.Background(), user.ID)
}

func (app *application) recordFailedLogin(user *database.User) error {
	lockedUntil := time.Now().Add(LockedUntil)
	u, err := app.queries.RecordFailedLogin(
		context.Background(),
		database.RecordFailedLoginParams{
			UserID:      user.ID,
			LockedUntil: &lockedUntil,
			MaxAttempts: MaxFailedAttempts,
		},
	)
	if err != nil {
		return err
	}
	return newWrongPassErr(u.FailedAttempts)
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

func (app *application) enableMfa(secret string) error {
	return app.queries.EnableMFA(
		context.Background(),
		database.EnableMFAParams{ID: app.currentUser.id, TotpSecret: &secret},
	)
}

func (app *application) disableMfa() error {
	return app.queries.DisableMFA(
		context.Background(),
		app.currentUser.id,
	)
}
