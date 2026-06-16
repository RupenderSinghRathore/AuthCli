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
		username, password, err := app.getUserPass()
		if err != nil {
			return "error reading username or password", err
		}
		user, err := app.register(username, password)
		if err != nil {
			return "failed to register", err
		}

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
