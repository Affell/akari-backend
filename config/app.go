package config

import (
	"os"
	"time"
)

type App struct {
	Port        string
	DebugLevel  string
	StartedTime time.Time
}

func InitApp() (app App) {
	app.StartedTime = time.Now()

	if env := os.Getenv("APP_LOG_LEVEL"); env != "" {
		app.DebugLevel = env
	} else {
		app.DebugLevel = "debug"
	}

	if env := os.Getenv("APP_PORT"); env != "" {
		app.Port = env
	} else {
		app.Port = "4000"
	}

	return
}
