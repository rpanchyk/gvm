package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/services"
	"github.com/rpanchyk/gvm/internal/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install specified Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installer := &services.Installer{Config: &utils.Config}
		if err := installer.Install(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
