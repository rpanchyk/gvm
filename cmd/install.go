package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/clients"
	"github.com/rpanchyk/gvm/internal/services/cacher"
	"github.com/rpanchyk/gvm/internal/services/downloader"
	"github.com/rpanchyk/gvm/internal/services/installer"
	"github.com/rpanchyk/gvm/internal/services/lister"
	"github.com/rpanchyk/gvm/internal/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install specified Go version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installer := installer.NewDefaultInstaller(
			&utils.Config,
			downloader.NewDefaultDownloader(
				&utils.Config,
				lister.NewDefaultListFetcher(
					&utils.Config,
					&clients.SimpleHttpClient{},
					cacher.NewDefaultListCacher(&utils.Config)),
				&clients.SimpleHttpSaver{}),
		)
		if err := installer.Install(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
