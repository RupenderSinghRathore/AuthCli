package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"RupenderSinghRathore/AuthCli/internal/database"

	"github.com/chzyer/readline"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	AppName           = "AuthCLI"
	MaxFailedAttempts = 3
	LockedUntil       = 7 * 24 * time.Hour
	SessionValidPeriod = 7 * 24 * time.Hour
)

type application struct {
	readWriter *readline.Instance
	quit       bool
	cfg        *configuration
	queries     *database.Queries
	currentUser       struct {
		id         int64
		name       string
		isLoggedIn bool
		mfaEnabled bool
	}
}

type configuration struct {
	db struct {
		dsn         string
		maxOpenConn int
		maxIdleTime time.Duration
	}
}

func main() {
	godotenv.Load()

	var cfg configuration

	dsn, err := getEnv("DATABASE_PATH")
	if err != nil {
		log.Fatal(err)
	}
	cfg.db.dsn = dsn

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
	app.cfg = &cfg
	app.queries = database.New(db)

	user, err := app.getSessionUser()
	switch {
	case err == nil:
		app.setCurrentUser(user)
		fmt.Printf("Welcome back %s to %s\n", app.currentUser.name, AppName)
	case errors.Is(err, ErrUserNotFound), errors.Is(err, os.ErrNotExist):
		fmt.Printf("Welcome to %s new user\n", AppName)
	default:
		log.Fatal(err)
	}

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

func openDB(cfg *configuration) (*sql.DB, error) {
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

func (app *application) setCurrentUser(user *database.User) {
	app.currentUser.isLoggedIn = true
	app.currentUser.name = user.Username
	app.currentUser.id = user.ID
	app.currentUser.mfaEnabled = user.MfaEnabled > 0
}

func (app *application) clearCurrentUser() {
	app.currentUser.isLoggedIn = false
	app.currentUser.name = ""
	app.currentUser.id = 0
	app.currentUser.mfaEnabled = false
}

func getEnv(v string) (string, error) {
	env := os.Getenv(v)
	if env == "" {
		return "", errors.New("Err: dsn environment variable not found")
	}
	return env, nil
}
