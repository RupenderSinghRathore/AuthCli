package main

import (
	"fmt"
)

func (app *application) execCmd(cmd string) (string, error) {
	var msg string
	var err error

	switch cmd {
	case "":

	case "exit":
		app.quit = true
		msg = "Bye.."

	case "register":
		return app.register()

	case "login":
	case "logout":
	case "whoami":
		msg = app.whoami()
	case "enable-2fa":
	case "disable-2fa":
	case "help":
	default:
		msg = fmt.Sprintf("no such commad: %s", cmd)
	}

	return msg, err
}

func (app *application) whoami() string {
	if app.user.isLoggedIn {
		return app.user.name
	} else {
		return "not logged in"
	}
}

func (app *application) register() (string, error) {
	if app.user.isLoggedIn {
		return "already registered", nil
	}
	username, password, err := app.getUserPass()
	if err != nil {
		return "", err
	}
	user, err := app.createUser(username, password)
	if err != nil {
		return "", err
	}
	app.loggingIn(user.Username)

	sessionId, err := app.createSession(user.ID)
	if err != nil {
		return "", err
	}

	if err := app.writeSessionConfig(sessionId); err != nil {
		return "", err
	}
	return "registered successfully!", nil
}
