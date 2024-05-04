package models

type Sdk struct {
	URL          string
	FilePath     string
	Version      string
	Os           string
	Arch         string
	IsDownloaded bool
	IsInstalled  bool
	IsDefault    bool
}
