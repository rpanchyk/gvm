package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "Go version manager",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot get current directory, error:", err)
		os.Exit(1)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Cannot get user home directory, error:", err)
		os.Exit(1)
	}

	viper.AddConfigPath(currentDir)
	viper.AddConfigPath(filepath.Join(userHomeDir, ".gvm"))
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Cannot read config, error:", err)
		os.Exit(1)
	}
}
