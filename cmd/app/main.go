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
