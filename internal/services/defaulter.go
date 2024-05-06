package services

import (
	"fmt"

	"github.com/rpanchyk/gvm/internal/models"
)

type Defaulter struct {
	Config *models.Config
}

func (d Defaulter) Default(version string) error {
	listFetcher := &ListFetcher{Config: d.Config}
	sdks, err := listFetcher.FetchAll()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}
	for _, sdk := range sdks {
		if sdk.Version == version && sdk.IsDefault {
			fmt.Printf("SDK version %s is already used as default\n", version)
			return nil
		}
	}

	platformDefaulter := &PlatformDefaulter{Config: d.Config}
	return platformDefaulter.Set(version)
}
