package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ildarusmanov/go-up/app"
)

type App struct {
	app.Application
}

func NewApp(ctx context.Context) *App {
	a := *app.NewApplication(ctx, nil, nil)

	a.SetConfig("userName", "User")
	a.SetConfig("appName", "Test App")

	a.AddServiceFactory("printer", PrinterFactory)
	a.AddServiceFactory("hello", HelloFactory)
	a.AddServiceFactory("server", ServerFactory)

	return &App{Application: a}
}

func (a *App) Stop() {
	srv, err := a.GetService("server")

	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.(*http.Server).Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("[*] Server stopped")
}
