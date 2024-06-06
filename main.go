package main

import (
	"akari/auth"
	"akari/config"
	"akari/handlers"
	"akari/handlers/wsHandler"
	"akari/models/grid"
	"akari/models/ws"
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

//go:embed public/*
var Folder embed.FS

func main() {
	// Lancement des routines
	golog.Info("launching routines...")

	// creation du nouveau Engine Iris
	router := iris.New()

	// DEBUG : utilisation du iris.Logger par défaut sinon rien
	router.Logger().SetLevel(config.Cfg.App.DebugLevel)
	if config.Cfg.App.DebugLevel == "debug" {
		router.Use(debugLogger())
	}

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "DELETE"},
		AllowCredentials: true,
	})
	router.UseRouter(crs)
	router.AllowMethods(iris.MethodOptions)
	router.Use(iris.Compression)

	router.Get("/favicon.ico", iris.Cache(30*time.Second), handlers.FaviconHandler(Folder))

	router.Get("/", iris.Cache(15*time.Second), handlers.IndexHandler)

	router.HandleDir("/img", handlers.FileHandler(Folder, "public/img"), iris.DirOptions{IndexName: "/img", Compress: true})

	wsRouter := ws.NewRouter()
	{
		wsRouter.On("auth", wsHandler.OnAuth, false)
		wsRouter.On("search", wsHandler.OnSearch, true)
		wsRouter.On("cancelSearch", wsHandler.OnCancelSearch, true)
		wsRouter.On("gridSubmit", wsHandler.OnGridSubmit, true)
		wsRouter.On("scoreboard", wsHandler.OnScoreboard, true)
		wsRouter.On("forfeit", wsHandler.OnForfeit, true)
	}
	router.Any("/ws", wsRouter.ServeHTTP)

	router.Use(auth.AuthRequired())
	{
		router.Any("/{service}", handlers.HandleRequest)
		router.Any("/{service}/{primary}", handlers.HandleRequest)
		router.Any("/{service}/{primary}/{secondary}", handlers.HandleRequest)
		router.Any("/{service}/{primary}/{secondary}/{tertiary}", handlers.HandleRequest)
		router.Any("/{service}/{primary}/{secondary}/{tertiary}/{tail}", handlers.HandleRequest)
	}

	// err := router.Listen(":" + config.Cfg.App.Port)
	// if err != nil {
	// 	golog.Fatal(err)
	// }

	g := grid.GenerateGrid(10, 1)
	golog.Debug(g)

	// var g grid.Grid
	// g = [][]int{
	// 	{-2, -2, -2, 2, -2, -2, 0, -1, -2, 1},
	// 	{-2, 1, 1, -2, -2, 1, 0, -2, -2, -2},
	// 	{-2, 1, -2, -2, 2, -2, 2, -2, -2, -2},
	// 	{1, -2, -2, -2, 0, 1, -2, 1, -2, -2},
	// 	{-2, -2, -2, -2, -2, -2, -2, -2, -2, -2},
	// 	{-2, -2, -2, -2, 0, 0, -2, -2, -2, -2},
	// 	{-2, -1, -2, 1, -2, -2, -2, -2, -2, -2},
	// 	{-2, -2, 1, -2, -2, 0, -2, 1, 0, 0},
	// 	{-2, -2, -2, 1, -2, -2, -2, -2, -2, 0},
	// 	{-2, -2, -2, 1, -2, -1, -2, -2, 0, -1},
	// }

	// var s grid.Grid
	// s = [][]int{
	// 	{-2, -2, -3, 2, -3, -2, 0, -1, -2, 1},
	// 	{-3, 1, 1, -2, -2, 1, 0, -2, -2, -3},
	// 	{-2, 1, -2, -3, 2, -3, 2, -3, -2, -2},
	// 	{1, -3, -2, -2, 0, 1, -2, 1, -2, -2},
	// 	{-2, -2, -2, -2, -2, -2, -2, -2, -3, -2},
	// 	{-2, -2, -3, -2, 0, 0, -2, -3, -2, -2},
	// 	{-2, -1, -2, 1, -2, -2, -3, -2, -2, -2},
	// 	{-2, -2, 1, -3, -2, 0, -2, 1, 0, 0},
	// 	{-2, -3, -2, 1, -2, -2, -2, -3, -2, 0},
	// 	{-3, -2, -2, 1, -3, -1, -2, -2, 0, -1},
	// }

	// golog.Debug(g.CheckSolution(s))

}

func debugLogger() iris.Handler {
	return func(c iris.Context) {
		t := time.Now()
		c.Next()
		params := []string{
			fmt.Sprint(c.GetStatusCode()),
			c.Request().Method,
			c.RequestPath(false),
			time.Since(t).String(),
		}
		golog.Debug(strings.Join(params, " | "))
	}
}
