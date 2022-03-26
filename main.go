package main

import (
	"github.com/JulesMike/spoty/bootstrap"
	"go.uber.org/fx"
)

// @title Spoty API
// @version 0.1.0
// @description Access information about current playing track on spotify through REST endpoints.

// @contact.name Jules Michael

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	fx.New(bootstrap.Module).Run()
}
