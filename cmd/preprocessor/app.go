package main

import (
	"agregator/preprocessor/internal/cmd/app"
	"time"
)

func main() {
	app, err := app.New(20 * time.Second)
	if err != nil {
		panic(err)
	}
	app.Run()
}
