package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: loadRoutes(),
	}
	return app
}

// Receiver/ Owner of the method similar to this keyword
func (a *App) Start(ctx context.Context) error {

	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to listen to server", err)
	}
	return nil
}
