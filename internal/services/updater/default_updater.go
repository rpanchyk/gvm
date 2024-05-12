package updater

import (
	"fmt"

	"github.com/rpanchyk/gvm/internal/models"
	"github.com/rpanchyk/gvm/internal/services/defaulter"
	"github.com/rpanchyk/gvm/internal/services/installer"
	"github.com/rpanchyk/gvm/internal/services/lister"
)

type DefaultUpdater struct {
	config      *models.Config
	listFetcher lister.ListFetcher
	installer   installer.Installer
	defaulter   defaulter.Defaulter
}

func NewDefaultUpdater(
	config *models.Config,
	listFetcher lister.ListFetcher,
	installer installer.Installer,
	defaulter defaulter.Defaulter) *DefaultUpdater {

	return &DefaultUpdater{
		config:      config,
		listFetcher: listFetcher,
		installer:   installer,
		defaulter:   defaulter,
	}
}

func (u DefaultUpdater) Update() error {
	sdks, err := u.listFetcher.Fetch()
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

	if err = u.installer.Install(version); err != nil {
		return fmt.Errorf("cannot install SDK %s version: %w", version, err)
	}

	if err = u.defaulter.Default(version); err != nil {
		return fmt.Errorf("cannot set default SDK %s version: %w", version, err)
	}

	return nil
}
