package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version of gvm",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gvm version %s %s/%s\n", "1.0.0", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
