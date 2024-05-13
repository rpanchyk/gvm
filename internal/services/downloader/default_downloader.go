package downloader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/rpanchyk/gvm/internal/clients"
	"github.com/rpanchyk/gvm/internal/models"
	"github.com/rpanchyk/gvm/internal/services/lister"
)

type DefaultDownloader struct {
	config      *models.Config
	listFetcher lister.ListFetcher
	httpSaver   clients.HttpSaver
}

func NewDefaultDownloader(
	config *models.Config,
	listFetcher lister.ListFetcher,
	httpSaver clients.HttpSaver) *DefaultDownloader {

	return &DefaultDownloader{
		config:      config,
		listFetcher: listFetcher,
		httpSaver:   httpSaver,
	}
}

func (d DefaultDownloader) Download(version string) (*models.Sdk, error) {
	sdks, err := d.listFetcher.Fetch()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	sdk, err := d.findSdk(version, sdks)
	if err != nil {
		return nil, fmt.Errorf("cannot find specified SDK: %w", err)
	}
	fmt.Printf("Found SDK: %+v\n", *sdk)

	filePath, err := d.downloadSdk(sdk.URL, d.config.DownloadDir)
	if err != nil {
		return nil, fmt.Errorf("cannot download specified SDK: %w", err)
	}

	sdk.FilePath = filePath
	sdk.IsDownloaded = true

	fmt.Printf("Downloaded SDK: %+v\n", *sdk)
	return sdk, nil
}

func (d DefaultDownloader) findSdk(version string, sdks []models.Sdk) (*models.Sdk, error) {
	for _, sdk := range sdks {
		if sdk.Version == version {
			return &sdk, nil
		}
	}
	return nil, fmt.Errorf("version %s not found", version)
}

func (d DefaultDownloader) downloadSdk(url, dir string) (string, error) {
	filePath := filepath.Join(dir, path.Base(url))
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("SDK %s has been already downloaded\n", filePath)
		return filePath, nil
	}

	if err := d.httpSaver.Save(url, filePath); err != nil {
		return "", fmt.Errorf("cannot save file to: %s error: %w", filePath, err)
	}

	fmt.Printf("SDK %s has been downloaded\n", filePath)
	return filePath, nil
}
