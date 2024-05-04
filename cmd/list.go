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
			os.Exit(1)
		}
		fmt.Printf("Parsed config: %+v\n", *config)

		listFetcher := services.ListFetcher{Config: config}
		sdks, err := listFetcher.Fetch()
		if err != nil {
			os.Exit(1)
		}

		for _, sdk := range sdks {
			fmt.Println(" ", sdk.Version, sdk.Os, sdk.Arch, sdk.URL)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
