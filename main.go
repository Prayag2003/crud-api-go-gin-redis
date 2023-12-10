package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	application "github.com/Prayag2003/crud-api-golang/Application"
)

func main() {
	app := application.New()

	ctx := context.Background()
	newCtx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err := app.Start(newCtx)
	if err != nil {
		fmt.Print("failed to start app %w", err)
	}
}
