package models

// Go SDK model
type Sdk struct {
	URL          string `json:"url"`
	Version      string `json:"version"`
	Os           string `json:"os"`
	Arch         string `json:"arch"`
	FilePath     string `json:"-"`
	IsDownloaded bool   `json:"-"`
	IsInstalled  bool   `json:"-"`
	IsDefault    bool   `json:"-"`
}
