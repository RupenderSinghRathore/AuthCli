package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"RupenderSinghRathore/AuthCli/internal/database"

	"github.com/chzyer/readline"
	_ "github.com/mattn/go-sqlite3"
)

const AppName = "authcli"

type application struct {
	readWriter *readline.Instance
	quit       bool
	cfg        *confugration
	queary     *database.Queries
}

type confugration struct {
	db struct {
		dsn         string
		maxOpenConn int
		maxIdleTime time.Duration
	}
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

func (app *application) getUserPass() (string, []byte, error) {
	oldPrompt := app.readWriter.Config.Prompt
	app.readWriter.SetPrompt("login: ")
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

func main() {
	fmt.Printf("Wellcome to %s\n", AppName)

	var cfg confugration

	flag.IntVar(&cfg.db.maxOpenConn, "db-max-open-conns", 25, "Sqlite max open connections")
	flag.DurationVar(
		&cfg.db.maxIdleTime,
		"db-max-idle-time",
		15*time.Minute,
		"Sqlite max connection idle time",
	)

	db, err := openDB(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := application{}
	app.queary = database.New(db)

	if err := app.repl(); err != nil {
		switch {
		case errors.Is(err, io.EOF):
		case errors.Is(err, readline.ErrInterrupt):
			os.Exit(1)
		default:
			log.Fatal(err)
		}
	}
}

func openDB(cfg *confugration) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConn)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	return db, nil
}
