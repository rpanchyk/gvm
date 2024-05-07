package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rpanchyk/gvm/internal/models"
	"github.com/rpanchyk/gvm/internal/utils"
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

	utils.Config = getConfig()
	fmt.Printf("Config: %+v\n", utils.Config)
}

func getConfig() models.Config {
	var config models.Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Cannot unmarshal config, error:", err)
		os.Exit(1)
	}

	config.DownloadDir = toAbsPath(config.DownloadDir)
	config.InstallDir = toAbsPath(config.InstallDir)
	config.LocalDir = toAbsPath(config.LocalDir)
	return config
}

func toAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		configFile := viper.GetViper().ConfigFileUsed()
		return filepath.Join(filepath.Dir(configFile), path)
	}
	return path
}
