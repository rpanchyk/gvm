package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/rpanchyk/gvm/internal/models"
)

type Downloader struct {
	Config *models.Config
}

func (d Downloader) Download(version string) (*models.Sdk, error) {
	listFetcher := &ListFetcher{Config: d.Config}
	sdks, err := listFetcher.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of SDKs: %w", err)
	}

	sdk, err := d.findSdk(version, sdks)
	if err != nil {
		return nil, fmt.Errorf("cannot find specified SDK: %w", err)
	}
	fmt.Printf("Found SDK: %+v\n", *sdk)

	filePath, err := d.downloadSdk(sdk.URL, d.Config.DownloadDir)
	if err != nil {
		return nil, fmt.Errorf("cannot download specified SDK: %w", err)
	}

	sdk.FilePath = filePath
	sdk.IsDownloaded = true

	fmt.Printf("Downloaded SDK: %+v\n", *sdk)
	return sdk, nil
}

func (d Downloader) findSdk(version string, sdks []models.Sdk) (*models.Sdk, error) {
	for _, sdk := range sdks {
		if sdk.Version == version && sdk.Os == runtime.GOOS && sdk.Arch == runtime.GOARCH {
			return &sdk, nil
		}
	}
	return nil, fmt.Errorf("version %s not found", version)
}

func (d Downloader) downloadSdk(fileUrl, dir string) (string, error) {
	fileName := path.Base(fileUrl)
	filePath := filepath.Join(dir, fileName)
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("SDK %s has been already downloaded\n", filePath)
		return filePath, nil
	}

	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", fmt.Errorf("cannot get data from url: %s error: %w", fileUrl, err)
	}
	defer resp.Body.Close()

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create dir: %s error: %w", dir, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %s error: %w", filePath, err)
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("cannot save file: %s error: %w", filePath, err)
	}

	fmt.Printf("SDK %s has been downloaded\n", filePath)
	return filePath, nil
}