package models

type Config struct {
	ReleaseURL  string `mapstructure:"release_url"`
	DownloadDir string `mapstructure:"download_dir"`
	InstallDir  string `mapstructure:"install_dir"`
	LocalDir    string `mapstructure:"local_dir"`
	Limit       int    `mapstructure:"limit"`
	FilterOs    bool   `mapstructure:"filter_os"`
	FilterArch  bool   `mapstructure:"filter_arch"`
}
