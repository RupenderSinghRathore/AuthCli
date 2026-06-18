package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func (app *application) execCommand(cmd string) (string, error) {
	var msg string
	var err error

	switch cmd {
	case "":

	case "exit":
		return app.exit(), nil

	case "register":
		return app.register()

	case "login":
		return app.login()

	case "logout":
		return app.logout()

	case "whoami":
		return app.whoami()

	case "enable-2fa":
		return app.enable2fa()

	case "disable-2fa":

	case "help":
		return app.help(), nil

	default:
		msg = fmt.Sprintf("no such commad: %s", cmd)
	}

	return msg, err
}

func (app *application) exit() string {
	app.quit = true
	return "bye.."
}

func (app *application) enable2fa() (string, error) {
	if !app.user.isLoggedIn {
		return "", ErrNotLoggedIn
	}
	if app.user.mfaEnabled {
		return "already enabled", nil
	}
	secret, err := generateTOTP(app.user.name)
	if err != nil {
		return "", err
	}

	// output the secert and info to start 2fa
	app.revealTotp(secret)
	code, err := app.promptTotp()
	if err != nil {
		return "", err
	}
	if !verifyTOTP(secret, code) {
		return "", ErrIncorrectCode
	}

	// store the secret in database
	if err := app.enableMfa(secret); err != nil {
		return "", err
	}
	app.user.mfaEnabled = true

	return "enabled TOTP based 2FA", nil
}

func (app *application) disable2fa() (string, error) {
	if !app.user.mfaEnabled {
		return "2FA not enabled", nil
	}

	err := app.disableMfa()
	if err != nil {
		return "", err
	}
	app.user.mfaEnabled = false

	return "disabled TOTP based 2FA", nil
}

func (app *application) whoami() (string, error) {
	if app.user.isLoggedIn {
		msg := "username: " + app.user.name + "\n"
		if app.user.mfaEnabled {
			msg += "2FA: enabled"
		} else {
			msg += "2FA: disabled"
		}
		return msg, nil
	} else {
		return "", ErrNotLoggedIn
	}
}

func (app *application) register() (string, error) {
	if app.user.isLoggedIn {
		return "already logged in", nil
	}
	username, password, err := app.promptUserPass()
	if err != nil {
		return "", err
	}

	if err := validateUsernamePassowrd(username, password); err != nil {
		return "", err
	}

	user, err := app.createUser(username, password)
	switch {
	case isUniqueConstraintErr(err):
		return "", ErrUsernameAlreadyExists
	case err != nil:
		return "", err
	}

	app.fillLoginInfo(user)

	sessionId, err := app.createSession(user.ID)
	if err != nil {
		return "", err
	}

	if err := app.writeSessionConfig(sessionId); err != nil {
		return "", err
	}
	return "registered successfully!", nil
}

func isUniqueConstraintErr(err error) bool {
	sqliteErr, ok := errors.AsType[sqlite3.Error](err)
	return ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}

func (app *application) login() (string, error) {
	if app.user.isLoggedIn {
		return "already logged in", nil
	}
	username, password, err := app.promptUserPass()
	if err != nil {
		return "", err
	}

	if err := validateUsernamePassowrd(username, password); err != nil {
		return "", err
	}

	user, err := app.getUserForLogin(username)
	if err != nil {
		return "", err
	}

	if !VerifyPassword(user.PasswordHash, string(password)) {
		return "", app.recordFailedLogin(user)
	}

	if user.MfaEnabled > 0 {
		code, err := app.promptTotp()
		if err != nil {
			return "", err
		}
		if !verifyTOTP(*user.TotpSecret, code) {
			app.recordFailedLogin(user)
			return "", ErrIncorrectCode
		}
	}

	app.fillLoginInfo(user)

	sessionId, err := app.createSession(user.ID)
	if err != nil {
		return "", err
	}

	if err := app.writeSessionConfig(sessionId); err != nil {
		return "", err
	}
	return "logged in as " + app.user.name, nil
}

func (app *application) logout() (string, error) {
	if !app.user.isLoggedIn {
		return "", ErrNotLoggedIn
	}
	if err := app.queary.DeleteSession(context.Background(), app.user.id); err != nil {
		return "", err
	}
	app.unFillLoginInfo()
	return "logged out successfully", nil
}

func (app *application) help() string {
	if app.user.isLoggedIn {
		return `
Available Commands

Account:
  whoami        Show current user details
  enable-2fa    Enable TOTP-based MFA
  disable-2fa   Disable MFA
  logout        End current session
  help          Show this help message
  exit          Quit the program
`
	}
	return `
Available Commands

Authentication:
  register      Create a new user account
  login         Login with username and password
  help          Show this help message
  exit          Quit the program
`
}
