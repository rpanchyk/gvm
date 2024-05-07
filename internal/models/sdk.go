package models

type Sdk struct {
	URL          string `json:"url"`
	FilePath     string `json:"-"`
	Version      string `json:"version"`
	Os           string `json:"os"`
	Arch         string `json:"arch"`
	IsDownloaded bool   `json:"-"`
	IsInstalled  bool   `json:"-"`
	IsDefault    bool   `json:"-"`
}
