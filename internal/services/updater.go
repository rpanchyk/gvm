package services

import (
	"fmt"

	"github.com/rpanchyk/gvm/internal/models"
)

type Updater struct {
	Config *models.Config
}

func (u Updater) Update() error {
	listFetcher := &ListFetcher{Config: u.Config}
	sdks, err := listFetcher.Fetch()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	version := sdks[0].Version
	for _, sdk := range sdks {
		if sdk.IsDefault && sdk.Version == version {
			fmt.Printf("SDK version %s is already the latest\n", version)
			return nil
		}
	}

	installer := &Installer{Config: u.Config}
	if err = installer.Install(version); err != nil {
		return fmt.Errorf("cannot install SDK %s version: %w", version, err)
	}

	defaulter := &Defaulter{Config: u.Config}
	if err = defaulter.Default(version); err != nil {
		return fmt.Errorf("cannot set default SDK %s version: %w", version, err)
	}

	return nil
}
