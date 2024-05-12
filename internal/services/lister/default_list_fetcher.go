package lister

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"

	"github.com/rpanchyk/gvm/internal/clients"
	"github.com/rpanchyk/gvm/internal/models"
	"github.com/rpanchyk/gvm/internal/services/cacher"
)

type DefaultListFetcher struct {
	config     *models.Config
	httpClient clients.HttpClient
	listCacher cacher.ListCacher
}

func NewDefaultListFetcher(
	config *models.Config,
	httpClient clients.HttpClient,
	listCacher cacher.ListCacher) *DefaultListFetcher {

	return &DefaultListFetcher{
		config:     config,
		httpClient: httpClient,
		listCacher: listCacher,
	}
}

func (f DefaultListFetcher) Fetch() ([]models.Sdk, error) {
	sdks, err := f.listCacher.Get()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of SDKs from cache: %w", err)
	}

	if len(sdks) == 0 {
		sdks, err = f.downloadSdks()
		if err != nil {
			return nil, fmt.Errorf("cannot download list of SDKs: %w", err)
		}
		if err := f.listCacher.Save(sdks); err != nil {
			return nil, fmt.Errorf("cannot save list of SDKs to cache: %w", err)
		}
	}

	return f.enrichSdks(sdks)
}

func (f DefaultListFetcher) downloadSdks() ([]models.Sdk, error) {
	response, err := f.httpClient.Get(f.config.ReleaseURL)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}

	sdks, err := f.parsePage(response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return sdks, nil
}

func (f DefaultListFetcher) parsePage(body string) ([]models.Sdk, error) {
	r, err := regexp.Compile(`href=['"]\/dl(/go([0-9.]*?)\.(\w+)-([\w\-.]+)\.(?:tar\.gz|zip)+)['"]`)
	if err != nil {
		return nil, fmt.Errorf("error compile regexp: %w", err)
	}

	sdks := make([]models.Sdk, 0)
	for _, parts := range r.FindAllStringSubmatch(body, -1) {
		// fmt.Printf("%v\n", parts)

		url, err := url.JoinPath(f.config.ReleaseURL, parts[1])
		if err != nil {
			return nil, fmt.Errorf("error composing url: %w", err)
		}

		sdk := models.Sdk{
			URL:     url,
			Version: parts[2],
			Os:      parts[3],
			Arch:    parts[4],
		}

		if f.config.ListFilterOs && runtime.GOOS != sdk.Os {
			continue
		}
		if f.config.ListFilterArch && runtime.GOARCH != sdk.Arch {
			continue
		}

		if !slices.Contains(sdks, sdk) {
			sdks = append(sdks, sdk)
		}
	}

	return sdks, nil
}

func (f DefaultListFetcher) enrichSdks(sdks []models.Sdk) ([]models.Sdk, error) {
	for i := 0; i < len(sdks); i++ {
		downloadedFile := filepath.Join(f.config.DownloadDir, filepath.Base(sdks[i].URL))
		if _, err := os.Stat(downloadedFile); err == nil {
			sdks[i].IsDownloaded = true
		}
		sdks[i].FilePath = downloadedFile

		goRootDir := filepath.Join(f.config.InstallDir, "go"+sdks[i].Version)
		if _, err := os.Stat(goRootDir); err == nil {
			sdks[i].IsInstalled = true
		}

		if envVar, ok := os.LookupEnv("GOROOT"); ok {
			if envVar == goRootDir {
				sdks[i].IsDefault = true
			} else {
				goPathDir := filepath.Join(f.config.LocalDir, "go"+sdks[i].Version)
				if _, err := os.Stat(goPathDir); err == nil && strings.HasPrefix(envVar, goPathDir) {
					sdks[i].IsDefault = true
				}
			}
		}
	}
	return sdks, nil
}
