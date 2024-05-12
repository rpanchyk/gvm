package lister

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/rpanchyk/gvm/internal/models"
)

type FilteredListFetcher struct {
	config      *models.Config
	listFetcher ListFetcher
}

func NewFilteredListFetcher(
	config *models.Config,
	listFetcher ListFetcher) *FilteredListFetcher {

	return &FilteredListFetcher{
		config:      config,
		listFetcher: listFetcher,
	}
}

func (f FilteredListFetcher) Fetch() ([]models.Sdk, error) {
	sdks, err := f.listFetcher.Fetch()
	if err != nil {
		return nil, fmt.Errorf("error fetching SDKs: %w", err)
	}

	filtered, err := f.filterSdks(sdks)
	if err != nil {
		return nil, fmt.Errorf("error filtering SDKs: %w", err)
	}
	return filtered, nil
}

func (f FilteredListFetcher) filterSdks(sdks []models.Sdk) ([]models.Sdk, error) {
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
		if count > f.config.ListLimit {
			break
		}
		res = append(res, sdk)
		ver = sdk.Version
	}

	clear(sdks)
	return res, nil
}
