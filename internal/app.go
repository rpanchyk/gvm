package internal

import (
	"fmt"
	"path/filepath"

	"github.com/rpanchyk/gvm/internal/models"
	"github.com/spf13/viper"
)

type App struct {
}

func (a App) Run() error {
	println("run")

	var config models.Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	// Justify sdk dir path
	if !filepath.IsAbs(config.Main.SdkDir) {
		configFile := viper.GetViper().ConfigFileUsed()
		config.Main.SdkDir = filepath.Join(filepath.Dir(configFile), config.Main.SdkDir)
	}

	fmt.Printf("Parsed %#v\n", config)
	return nil
}
