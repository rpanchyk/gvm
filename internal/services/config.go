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

	// Justify sdk dir path
	if !filepath.IsAbs(config.SdkDir) {
		configFile := viper.GetViper().ConfigFileUsed()
		config.SdkDir = filepath.Join(filepath.Dir(configFile), config.SdkDir)
	}

	return &config, nil
}
