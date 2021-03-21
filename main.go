package main

import (
	"log"

	"github.com/JulesMike/spoty/app"
	"github.com/JulesMike/spoty/server"
	"github.com/JulesMike/spoty/spoty"
	"go.uber.org/fx"
)

// @title Spoty API
// @version 0.1.0
// @description Access information about current playing track on spotify through REST endpoints.

// @contact.name Jules Michael

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	if err := run(); err != nil {
		log.Fatalf("[Spoty] %v", err)
	}
}

func run() error {
	app := fx.New(
		app.DefaultProviders,
		fx.Invoke(func(spoty *spoty.Spoty, server *server.Server) {
			spoty.RegisterRoutes(server.APIRoute())
		}),
	)
	if err := app.Err(); err != nil {
		return err
	}

	app.Run()

	return nil
}
