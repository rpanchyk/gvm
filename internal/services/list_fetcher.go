package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/rpanchyk/gvm/internal/models"
)

type ListFetcher struct {
	Config *models.Config
}

func (f ListFetcher) Fetch() ([]models.Sdk, error) {
	sdks, err := f.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("error fetching all SDKs: %w", err)
	}

	reduced, err := f.reduce(sdks)
	if err != nil {
		return nil, fmt.Errorf("error filtering SDKs: %w", err)
	}
	return reduced, nil
}

func (f ListFetcher) FetchAll() ([]models.Sdk, error) {
	response, err := http.Get(f.Config.ReleaseURL)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %w", err)
	}
	defer response.Body.Close()
	// fmt.Printf("status code: %d\n", response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	// fmt.Printf("Body: %s\n", body)

	sdks, err := f.parse(string(body))
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	for i := 0; i < len(sdks); i++ {
		downloadedFile := filepath.Join(f.Config.DownloadDir, filepath.Base(sdks[i].URL))
		if _, err := os.Stat(downloadedFile); err == nil {
			sdks[i].IsDownloaded = true
		}
		sdks[i].FilePath = downloadedFile

		goRootDir := filepath.Join(f.Config.InstallDir, "go"+sdks[i].Version)
		if _, err := os.Stat(goRootDir); err == nil {
			sdks[i].IsInstalled = true
		}

		if envVar, ok := os.LookupEnv("GOROOT"); ok {
			if envVar == goRootDir {
				sdks[i].IsDefault = true
			} else {
				goPathDir := filepath.Join(f.Config.LocalDir, "go"+sdks[i].Version)
				if _, err := os.Stat(goPathDir); err == nil && strings.HasPrefix(envVar, goPathDir) {
					sdks[i].IsDefault = true
				}
			}
		}
	}

	return sdks, nil
}

func (f ListFetcher) parse(s string) ([]models.Sdk, error) {
	sdks := make([]models.Sdk, 0)

	r, err := regexp.Compile(`href=['"]\/dl(/go([0-9.]*?)\.(\w+)-([\w\-.]+)\.(?:tar\.gz|zip)+)['"]`)
	if err != nil {
		return nil, fmt.Errorf("error compile regexp: %w", err)
	}

	for _, parts := range r.FindAllStringSubmatch(s, -1) {
		// fmt.Printf("%v\n", parts)

		url, err := url.JoinPath(f.Config.ReleaseURL, parts[1])
		if err != nil {
			return nil, fmt.Errorf("error composing url: %w", err)
		}

		sdk := models.Sdk{
			URL:     url,
			Version: parts[2],
			Os:      parts[3],
			Arch:    parts[4],
		}

		if f.Config.FilterOs && runtime.GOOS != sdk.Os {
			continue
		}
		if f.Config.FilterArch && runtime.GOARCH != sdk.Arch {
			continue
		}

		sdks = append(sdks, sdk)
	}

	return sdks, nil
}

func (f ListFetcher) reduce(sdks []models.Sdk) ([]models.Sdk, error) {
	sort.Slice(sdks, func(i, j int) bool {
		first := strings.Split(sdks[i].Version, ".")
		second := strings.Split(sdks[j].Version, ".")

		length := max(len(first), len(second))
		for k := 0; k < length; k++ {
			if len(first) > k+1 && len(second) <= k+1 { // 1.9.1 vs 1.9
				return true
			}

			if len(first) <= k+1 && len(second) > k+1 { // 1.9 vs 1.9.1
				return false
			}

			if first[k] != second[k] { // 1.9.1 vs 1.9.2
				f, err := strconv.Atoi(first[k])
				if err != nil {
					panic(err)
				}
				s, err := strconv.Atoi(second[k])
				if err != nil {
					panic(err)
				}
				return f > s
			}
		}

		return false
	})
	// fmt.Printf("Sorted sdks: %v\n", sdks)

	res := make([]models.Sdk, 0)
	count := 0
	ver := ""
	for _, sdk := range sdks {
		if ver != sdk.Version {
			count++
		}
		if count > f.Config.Limit {
			break
		}
		res = append(res, sdk)
		ver = sdk.Version
	}

	clear(sdks)
	return res, nil
}
