package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version of gvm",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.0.0")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
