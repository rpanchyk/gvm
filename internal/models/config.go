package models

type Config struct {
	AllReleasesURL string `mapstructure:"all_releases_url"`
	DownloadDir    string `mapstructure:"download_dir"`
	InstallDir     string `mapstructure:"install_dir"`
	Limit          int    `mapstructure:"limit"`
	FilterOs       bool   `mapstructure:"filter_os"`
	FilterArch     bool   `mapstructure:"filter_arch"`
}
