package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/clients"
	"github.com/rpanchyk/gvm/internal/services/cacher"
	"github.com/rpanchyk/gvm/internal/services/lister"
	"github.com/rpanchyk/gvm/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows list of available Go versions",
	Run: func(cmd *cobra.Command, args []string) {
		listFetcher := lister.NewFilteredListFetcher(
			&utils.Config,
			lister.NewDefaultListFetcher(
				&utils.Config,
				&clients.SimpleHttpClient{},
				cacher.NewDefaultListCacher(&utils.Config)),
		)
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
			downloadedMarker := "            "
			if sdk.IsDownloaded {
				downloadedMarker = "[downloaded]"
			}
			installedMarker := ""
			if sdk.IsInstalled {
				installedMarker = "[installed]"
			}
			fmt.Printf("%s %s\t%s\t%s\t%s %s\n",
				defaultMarker, sdk.Version, sdk.Os, sdk.Arch, downloadedMarker, installedMarker)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
