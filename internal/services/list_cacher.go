package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rpanchyk/gvm/internal/models"
)

type ListCacher struct {
	CacheFile string
	TTL       time.Duration
}

func (c ListCacher) Get() ([]models.Sdk, error) {
	fileInfo, err := os.Stat(c.CacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Cache file not found: %s\n", c.CacheFile)
			return nil, nil
		}
		return nil, fmt.Errorf("cannot get cache file %s: %w", c.CacheFile, err)
	}
	if fileInfo.ModTime().Add(c.TTL).Before(time.Now()) { // expired
		fmt.Printf("Cache file found but expired: %s\n", c.CacheFile)
		return nil, nil
	}

	bytes, err := os.ReadFile(c.CacheFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read cache file %s: %w", c.CacheFile, err)
	}

	sdks := []models.Sdk{}
	if err := json.Unmarshal(bytes, &sdks); err != nil {
		return nil, fmt.Errorf("cannot decode cache from file %s: %w", c.CacheFile, err)
	}

	fmt.Printf("SDKs fetched from cache: %d\n", len(sdks))
	return sdks, nil
}

func (c ListCacher) Save(sdks []models.Sdk) error {
	file, err := os.Create(c.CacheFile)
	if err != nil {
		return fmt.Errorf("cannot create cache file: %w", err)
	}
	defer file.Close()

	b, err := json.MarshalIndent(sdks, "", "\t")
	if err != nil {
		return fmt.Errorf("cannot encode cache: %w", err)
	}

	r := bytes.NewReader(b)
	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("cannot persist cache: %w", err)
	}

	fmt.Printf("SDKs saved to cache: %d\n", len(sdks))
	return nil
}
