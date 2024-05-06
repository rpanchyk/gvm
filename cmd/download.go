package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/services"
	"github.com/rpanchyk/gvm/internal/utils"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download specified Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		downloader := &services.Downloader{Config: &utils.Config}
		if _, err := downloader.Download(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
}
