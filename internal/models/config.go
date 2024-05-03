package models

type Config struct {
	AllReleasesURL string `mapstructure:"all_releases_url"`
	SdkDir         string `mapstructure:"sdk_dir"`
	Limit          int    `mapstructure:"limit"`
	FilterOs       bool   `mapstructure:"filter_os"`
	FilterArch     bool   `mapstructure:"filter_arch"`
}
