package lister

import (
	"github.com/rpanchyk/gvm/internal/models"
)

type ListFetcher interface {
	Fetch() ([]models.Sdk, error)
}
