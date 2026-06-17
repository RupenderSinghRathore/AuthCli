package main

import (
	"RupenderSinghRathore/AuthCli/internal/database"
	"context"
	"os"
	"path/filepath"
)

func (app *application) writeSessionConfig(sessionId string) error {
	path, err := sessionPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(sessionId), 0o644)
}

func (app *application) readSessionConfig() ([]byte, error) {
	path, err := sessionPath()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func (app *application) createSession(userId int64) (string, error) {
	return app.queary.CreateSession(context.Background(), userId)
}

func sessionPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	sessionPath := filepath.Join(home, ".config", AppName, "session")
	return sessionPath, nil
}

func (app *application) getSessionUser() (*database.User, error) {
	session_id, err := app.readSessionConfig()
	if err != nil || len(session_id) == 0 {
		return nil, err
	}
	user, err := app.queary.GetUserBySessionToken(context.Background(), string(session_id))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
