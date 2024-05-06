package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/services"
	"github.com/rpanchyk/gvm/internal/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install the latest Go version",
	Run: func(cmd *cobra.Command, args []string) {
		updater := &services.Updater{Config: &utils.Config}
		if _, err := updater.Update(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
