package main

import (
	"agregator/preprocessor/internal/cmd/app"
	"log/slog"
	"time"
)

func main() {
	app, err := app.New(20*time.Second, slog.Default())
	if err != nil {
		panic(err)
	}
	app.Run()
}
