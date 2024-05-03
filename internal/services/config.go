package services

import (
	"path/filepath"

	"github.com/rpanchyk/gvm/internal/models"
	"github.com/spf13/viper"
)

type Config struct {
}

func (c Config) GetConfig() (*models.Config, error) {
	var config models.Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	if !filepath.IsAbs(config.DownloadDir) {
		configFile := viper.GetViper().ConfigFileUsed()
		config.DownloadDir = filepath.Join(filepath.Dir(configFile), config.DownloadDir)
	}
	if !filepath.IsAbs(config.InstallDir) {
		configFile := viper.GetViper().ConfigFileUsed()
		config.InstallDir = filepath.Join(filepath.Dir(configFile), config.InstallDir)
	}
	if !filepath.IsAbs(config.LocalDir) {
		configFile := viper.GetViper().ConfigFileUsed()
		config.LocalDir = filepath.Join(filepath.Dir(configFile), config.LocalDir)
	}

	return &config, nil
}
