package cmd

import (
	"fmt"
	"os"

	"github.com/JulesMike/spoty/build"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the build information",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := build.New()
		if err != nil {
			fmt.Print(err)
			os.Exit(1) //nolint:revive
		}

		fmt.Printf("Revision: %v\n", info.Revision)
		fmt.Printf("Last Commit: %v\n", info.LastCommit)
		fmt.Printf("Dirty Build: %v\n", info.DirtyBuild)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
