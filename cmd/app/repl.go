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

	var validErr ValidationErr
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
			errors.Is(err, ErrIncorrectCode),
			errors.Is(err, ErrUserNotFound),
			errors.Is(err, ErrUsernameAlreadyExists):
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

func (app *application) promptUserPass() (string, []byte, error) {
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

func (app *application) promptTotp() (string, error) {
	oldPrompt := app.readWriter.Config.Prompt
	app.readWriter.SetPrompt("your authenticator key: ")
	code, err := app.read()
	app.readWriter.SetPrompt(oldPrompt)
	return code, err
}

func (app *application) revealTotp(secret string) {
	app.write("use your authenticator app to setup totp with these values:")
	app.write(fmt.Sprintf("code name: %s", app.currentUser.name))
	app.write(fmt.Sprintf("your key: %s", secret))
}
