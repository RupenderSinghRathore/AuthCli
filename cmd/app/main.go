package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

const AppName = "authcli"

type application struct {
	readWriter *readline.Instance
	quit       bool
}

func sessionPath() (error, string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return err, ""
	}
	sessionPath := filepath.Join(home, ".config", AppName, "session")
	return nil, sessionPath
}

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

	for !app.quit {
		line, err := app.read()
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		msg, err := app.execCmd(line)
		if err != nil {
			return err
		}
		if msg != "" {
			app.write(fmt.Sprint(msg, "\n"))
		}
	}
	return nil
}

func (app *application) write(msg string) error {
	_, err := app.readWriter.Write([]byte(msg))
	return err
}

func (app *application) read() (string, error) {
	return app.readWriter.Readline()
}

func (app *application) execCmd(cmd string) (string, error) {
	var msg string
	var err error

	switch cmd {
	case "":

	case "exit":
		app.quit = true
		msg = "Bye.."

	case "register":
		// msg, err = app.promptUser("username: ")

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

func (app *application) promptUser(p string) (string, error) {
	oldPrompt := app.readWriter.Config.Prompt
	app.readWriter.SetPrompt(p)
	msg, err := app.read()
	app.readWriter.SetPrompt(oldPrompt)
	return msg, err
}

func main() {
	fmt.Printf("Wellcome to %s\n", AppName)
	app := application{}
	if err := app.repl(); err != nil {
		switch {
		case errors.Is(err, io.EOF):
		case errors.Is(err, readline.ErrInterrupt):
			os.Exit(1)
		}
		panic(err)
	}
}
