package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

func (app *application) repl() error {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("register"),
		readline.PcItem("login"),
		readline.PcItem("help"),
		readline.PcItem("exit"),

		readline.PcItem("whoami"),
		readline.PcItem("enable-2fa"),
		readline.PcItem("disable-2fa"),
		readline.PcItem("logout"),
	)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		AutoComplete:    completer,
		HistoryFile:     filepath.Join("/tmp", fmt.Sprint(AppName, ".history")),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return err
	}
	app.readWriter = rl
	defer rl.Close()

	var validErr ValidatationErr
	var lockErr AccountLockedErr
	var failedErr WrongPasswordErr

	for !app.quit {
		line, err := app.read()
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		msg, err := app.execCommand(line)
		switch {
		case err == nil:
			if msg != "" {
				app.write(msg)
			}
		case errors.As(err, &validErr),
			errors.As(err, &lockErr),
			errors.As(err, &failedErr),
			errors.Is(err, ErrNotLoggedIn),
			errors.Is(err, ErrUserNotFound):
			app.error(err)
		default:
			return err
		}
	}
	return nil
}

func (app *application) error(err error) {
	app.readWriter.Stderr().Write([]byte(err.Error() + "\n"))
}

func (app *application) write(msg string) error {
	_, err := app.readWriter.Write([]byte(msg + "\n"))
	return err
}

func (app *application) read() (string, error) {
	return app.readWriter.Readline()
}

func (app *application) getUserPass() (string, []byte, error) {
	oldPrompt := app.readWriter.Config.Prompt
	app.readWriter.SetPrompt("username: ")
	username, err := app.read()
	if err != nil {
		return "", nil, err
	}
	password, err := app.readWriter.ReadPassword("Password: ")
	if err != nil {
		return "", nil, err
	}
	app.readWriter.SetPrompt(oldPrompt)
	return username, password, err
}
