package main

import (
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
)

func main() {
	config := config.LoadConfig()
	app := app.App(config, zerolog.DebugLevel)
	app.Run()
}
