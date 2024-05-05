package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rpanchyk/gvm/internal/models"
)

type Remover struct {
	Config *models.Config
}

func (r Remover) Remove(removeDownloaded bool, removeInstalled bool, version string) error {
	listFetcher := &ListFetcher{Config: r.Config}
	sdks, err := listFetcher.FetchAll()
	if err != nil {
		return fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	var foundSdk *models.Sdk
	for _, sdk := range sdks {
		if sdk.Version == version {
			if sdk.IsDefault {
				return fmt.Errorf("cannot remove SDK version %s since it is used as default", version)
			}
			foundSdk = &sdk
			break
		}
	}

	if foundSdk == nil {
		return fmt.Errorf("cannot find SDK version %s", version)
	}
	fmt.Printf("Found SDK: %+v\n", *foundSdk)

	if removeDownloaded {
		if err := r.removeDownloaded(foundSdk); err != nil {
			return err
		}
	}

	if removeInstalled {
		if err := r.removeInstalled(foundSdk); err != nil {
			return err
		}
	}

	return nil
}

func (r Remover) removeDownloaded(sdk *models.Sdk) error {
	if !sdk.IsDownloaded {
		fmt.Printf("SDK %s version doesn't downloaded\n", sdk.Version)
		return nil
	}

	if err := os.Remove(sdk.FilePath); err != nil {
		return fmt.Errorf("cannot remove downloaded archive of SDK %s version: %w", sdk.Version, err)
	}

	fmt.Printf("Downloaded archive of SDK %s version has been removed\n", sdk.Version)
	return nil
}

func (r Remover) removeInstalled(sdk *models.Sdk) error {
	if !sdk.IsInstalled {
		fmt.Printf("SDK %s version doesn't installed\n", sdk.Version)
		return nil
	}

	goRootDir := filepath.Join(r.Config.InstallDir, "go"+sdk.Version)
	goPathDir := filepath.Join(r.Config.LocalDir, "go"+sdk.Version)
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
