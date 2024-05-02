package cmd

import (
	"os"

	"github.com/rpanchyk/gvm/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "Go version manager",
	Run: func(cmd *cobra.Command, args []string) {
		println("cmd start")

		app := &internal.App{}
		app.Run()

		println("cmd end")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
