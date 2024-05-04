package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/gvm/internal/services"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install the latest Go version",
	Run: func(cmd *cobra.Command, args []string) {
		configService := &services.Config{}
		config, err := configService.GetConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Parsed config: %+v\n", *config)

		updater := &services.Updater{Config: config}
		if _, err = updater.Update(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}