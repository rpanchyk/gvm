package models

type Sdk struct {
	URL          string
	Version      string
	Os           string
	Arch         string
	IsDownloaded bool
	IsInstalled  bool
}
