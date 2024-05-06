//go:build !windows

package services

import (
	"github.com/rpanchyk/gvm/internal/models"
)

type PlatformDefaulter struct {
	Config *models.Config
}

func (d PlatformDefaulter) Set(version string) error {
	panic("not implemented")
}
