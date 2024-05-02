package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rpanchyk/gvm/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "Go version manager",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()

		println("cmd start")

		app := &internal.App{}
		if err := app.Run(); err != nil {
			os.Exit(1)
		}

		println("cmd end")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Cannot get user home directory, error:", err)
		os.Exit(1)
	}

	viper.AddConfigPath(filepath.Join(userHomeDir, ".gvm"))
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Cannot read config, error:", err)
		os.Exit(1)
	}
}
