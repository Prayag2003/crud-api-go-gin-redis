package main

import (
	"context"
	"fmt"

	application "github.com/Prayag2003/go-microservice/Application"
)

func main() {
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("Failed to start app ")
	}
}
