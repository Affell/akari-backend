package main

import (
	"akari/config"
	"akari/models/postgresql"
	"time"

	"github.com/kataras/golog"
	"github.com/provectio/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		golog.Warn("No .env file to load")
	}

	config.Cfg.App = config.InitApp()
	golog.SetLevel(config.Cfg.App.DebugLevel)

	postgresql.SQLCtx, postgresql.SQLConn = config.InitPgSQL()

	golog.Debug("init success in " + time.Since(config.Cfg.App.StartedTime).String())
}
