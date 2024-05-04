package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/services"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows list of available Go versions",
	Run: func(cmd *cobra.Command, args []string) {
		configService := &services.Config{}
		config, err := configService.GetConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Parsed config: %+v\n", *config)

		listFetcher := services.ListFetcher{Config: config}
		sdks, err := listFetcher.Fetch()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, sdk := range sdks {
			defaultMarker := " "
			if sdk.IsDefault {
				defaultMarker = "*"
			}
			downloadedMarker := " "
			if sdk.IsDownloaded {
				downloadedMarker = "[downloaded]"
			}
			installedMarker := " "
			if sdk.IsInstalled {
				installedMarker = "[installed]"
			}
			fmt.Println(defaultMarker, sdk.Version, sdk.Os, sdk.Arch, downloadedMarker, installedMarker)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
