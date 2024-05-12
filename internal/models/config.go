package models

// Application config model
type Config struct {
	ReleaseURL      string `mapstructure:"release_url"`
	DownloadDir     string `mapstructure:"download_dir"`
	InstallDir      string `mapstructure:"install_dir"`
	LocalDir        string `mapstructure:"local_dir"`
	ListLimit       int    `mapstructure:"list_limit"`
	ListFilterOs    bool   `mapstructure:"list_filter_os"`
	ListFilterArch  bool   `mapstructure:"list_filter_arch"`
	ListCacheFile   string `mapstructure:"list_cache_file"`
	ListCacheTTL    int    `mapstructure:"list_cache_ttl"`
	UnixShellConfig string `mapstructure:"unix_shell_config"`
}
