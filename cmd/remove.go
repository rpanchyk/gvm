package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/services"
	"github.com/rpanchyk/gvm/internal/utils"
	"github.com/spf13/cobra"
)

var (
	removeDownloaded bool
	removeInstalled  bool
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove specified Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remover := &services.Remover{Config: &utils.Config}
		if err := remover.Remove(args[0], removeDownloaded, removeInstalled); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeDownloaded, "download", "d", false, "Remove downloaded SDK archive")
	removeCmd.Flags().BoolVarP(&removeInstalled, "install", "i", false, "Remove installed SDK directories")
	removeCmd.MarkFlagsOneRequired("download", "install")
	RootCmd.AddCommand(removeCmd)
}
