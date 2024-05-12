package downloader

import (
	"github.com/rpanchyk/gvm/internal/models"
)

type Downloader interface {
	Download(version string) (*models.Sdk, error)
}
