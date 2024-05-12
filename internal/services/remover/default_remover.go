package remover

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rpanchyk/gvm/internal/models"
	"github.com/rpanchyk/gvm/internal/services/lister"
)

type DefaultRemover struct {
	config      *models.Config
	listFetcher lister.ListFetcher
}

func NewDefaultRemover(
	config *models.Config,
	listFetcher lister.ListFetcher) *DefaultRemover {

	return &DefaultRemover{
		config:      config,
		listFetcher: listFetcher,
	}
}

func (r DefaultRemover) Remove(version string, removeDownloaded, removeInstalled bool) error {
	sdks, err := r.listFetcher.Fetch()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	sdk, err := r.findSdk(version, sdks)
	if err != nil {
		return fmt.Errorf("cannot find specified SDK: %w", err)
	}
	if sdk.IsDefault {
		return fmt.Errorf("cannot remove SDK version %s since it is used as default", version)
	}
	fmt.Printf("Found SDK: %+v\n", *sdk)

	if removeDownloaded {
		if err := r.removeDownloaded(sdk); err != nil {
			return err
		}
	}

	if removeInstalled {
		if err := r.removeInstalled(sdk); err != nil {
			return err
		}
	}

	return nil
}

func (r DefaultRemover) findSdk(version string, sdks []models.Sdk) (*models.Sdk, error) {
	for _, sdk := range sdks {
		if sdk.Version == version {
			return &sdk, nil
		}
	}
	return nil, fmt.Errorf("version %s not found", version)
}

func (r DefaultRemover) removeDownloaded(sdk *models.Sdk) error {
	if !sdk.IsDownloaded {
		fmt.Printf("SDK %s version is not downloaded\n", sdk.Version)
		return nil
	}

	if err := os.Remove(sdk.FilePath); err != nil {
		return fmt.Errorf("cannot remove downloaded archive of SDK %s version: %w", sdk.Version, err)
	}

	fmt.Printf("Downloaded archive of SDK %s version has been removed\n", sdk.Version)
	return nil
}

func (r DefaultRemover) removeInstalled(sdk *models.Sdk) error {
	if !sdk.IsInstalled {
		fmt.Printf("SDK %s version is not installed\n", sdk.Version)
		return nil
	}

	goRootDir := filepath.Join(r.config.InstallDir, "go"+sdk.Version)
	goPathDir := filepath.Join(r.config.LocalDir, "go"+sdk.Version)
	for _, dir := range []string{goRootDir, goPathDir} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("Directory %s doesn't exist\n", dir)
			continue
		}
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("cannot remove %s: %w", dir, err)
		}
		fmt.Printf("Directory %s has been removed\n", dir)
	}

	fmt.Printf("Installation directories of SDK %s version has been removed\n", sdk.Version)
	return nil
}
