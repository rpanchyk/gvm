package defaulter

import (
	"fmt"

	"github.com/rpanchyk/gvm/internal/models"
	"github.com/rpanchyk/gvm/internal/services/lister"
)

type DefaultDefaulter struct {
	config      *models.Config
	listFetcher lister.ListFetcher
}

func NewDefaultDefaulter(
	config *models.Config,
	listFetcher lister.ListFetcher) *DefaultDefaulter {

	return &DefaultDefaulter{
		config:      config,
		listFetcher: listFetcher,
	}
}

func (d DefaultDefaulter) Default(version string) error {
	sdks, err := d.listFetcher.Fetch()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	for _, sdk := range sdks {
		if sdk.Version == version && sdk.IsDefault {
			fmt.Printf("SDK version %s is already used as default\n", version)
			return nil
		}
	}

	platformDefaulter := &PlatformDefaulter{Config: d.config}
	return platformDefaulter.Default(version)
}
