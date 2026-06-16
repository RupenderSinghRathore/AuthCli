package main

import "fmt"

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
	case "enable-2fa":
	case "disable-2fa":
	case "help":
	default:
		msg = fmt.Sprintf("no such commad: %s", cmd)
	}

	return msg, err
}

func (app *application) register() (string, error) {
	username, password, err := app.getUserPass()
	if err != nil {
		return "", err
	}
	user, err := app.createUser(username, password)
	if err != nil {
		return "", err
	}

	// TODO: create session for the user
	_ = user

	return "registered successfully!", nil
}
