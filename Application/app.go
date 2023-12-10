package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{}),
	}

	app.loadRoutes()
	
	return app
}

// Receiver/ Owner of the method similar to this keyword
func (a *App) Start(ctx context.Context) error {

	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// NOTE: Checking if the Redis client is connected
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		fmt.Print("failed to start app %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis client")
		}
	}()

	fmt.Println("Server up and running at localhost:3000")

	// NOTE:
	ch := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to listen to server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch: // listener
		return err // return err to the caller
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(timeout)
	}

	return nil
}
