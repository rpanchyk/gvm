package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/clients"
	"github.com/rpanchyk/gvm/internal/services/cacher"
	"github.com/rpanchyk/gvm/internal/services/defaulter"
	"github.com/rpanchyk/gvm/internal/services/lister"
	"github.com/rpanchyk/gvm/internal/utils"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Set specified Go version as default",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defaulter := defaulter.NewDefaultDefaulter(
			&utils.Config,
			lister.NewDefaultListFetcher(
				&utils.Config,
				&clients.SimpleHttpClient{},
				cacher.NewDefaultListCacher(&utils.Config)),
		)
		if err := defaulter.Default(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(defaultCmd)
}
