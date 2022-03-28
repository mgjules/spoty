package cmd

import (
	"fmt"
	"os"

	"github.com/JulesMike/spoty/build"
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/logger"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the build information",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := logger.New(&config.Config{})
		if err != nil {
			fmt.Printf("failed to start logger: %v", err)
			os.Exit(1) //nolint:revive
		}

		info, err := build.New(logger)
		if err != nil {
			logger.Fatal(err)
		}

		fmt.Printf("Revision: %v\n", info.Revision)
		fmt.Printf("Last Commit: %v\n", info.LastCommit)
		fmt.Printf("Dirty Build: %v\n", info.DirtyBuild)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
