package main

import (
	"fmt"
	"os"

	"github.com/mgjules/spoty/cmd"
)

// @title Spoty API
// @version v0.2.1
// @description Access information about current playing track on spotify through REST endpoints.

// @contact.name Jules Michael

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("failed to execute cmd: %v", err)
		os.Exit(1)
	}
}
