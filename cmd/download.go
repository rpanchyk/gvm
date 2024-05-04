package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/services"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download specified Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configService := &services.Config{}
		config, err := configService.GetConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Parsed config: %+v\n", *config)

		downloader := &services.Downloader{Config: config}
		if err = downloader.Download(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
}
