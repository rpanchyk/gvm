package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/clients"
	"github.com/rpanchyk/gvm/internal/services/cacher"
	"github.com/rpanchyk/gvm/internal/services/defaulter"
	"github.com/rpanchyk/gvm/internal/services/downloader"
	"github.com/rpanchyk/gvm/internal/services/installer"
	"github.com/rpanchyk/gvm/internal/services/lister"
	"github.com/rpanchyk/gvm/internal/services/updater"
	"github.com/rpanchyk/gvm/internal/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Go to the latest version and set it as default",
	Run: func(cmd *cobra.Command, args []string) {
		updater := updater.NewDefaultUpdater(
			&utils.Config,

			lister.NewDefaultListFetcher(
				&utils.Config,
				&clients.SimpleHttpClient{},
				cacher.NewDefaultListCacher(&utils.Config)),

			installer.NewDefaultInstaller(
				&utils.Config,
				downloader.NewDefaultDownloader(
					&utils.Config,
					lister.NewDefaultListFetcher(
						&utils.Config,
						&clients.SimpleHttpClient{},
						cacher.NewDefaultListCacher(&utils.Config)),
					&clients.SimpleHttpSaver{})),

			defaulter.NewDefaultDefaulter(
				&utils.Config,
				lister.NewDefaultListFetcher(
					&utils.Config,
					&clients.SimpleHttpClient{},
					cacher.NewDefaultListCacher(&utils.Config))),
		)
		if err := updater.Update(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
