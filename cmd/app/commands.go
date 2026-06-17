package main

import (
	"context"
	"fmt"
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

func (app *application) whoami() (string, error) {
	if app.user.isLoggedIn {
		return app.user.name, nil
	} else {
		return "", ErrNotLoggedIn
	}
}

func (app *application) register() (string, error) {
	if app.user.isLoggedIn {
		return "already logged in", nil
	}
	username, password, err := app.getUserPass()
	if err != nil {
		return "", err
	}
	user, err := app.createUser(username, password)
	if err != nil {
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

func (app *application) login() (string, error) {
	if app.user.isLoggedIn {
		return "already logged in", nil
	}
	username, password, err := app.getUserPass()
	if err != nil {
		return "", err
	}

	user, err := app.tryUserLogin(username, password)
	if err != nil {
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
