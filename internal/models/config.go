package models

type Config struct {
	Main Main `mapstructure:"main"`
}

type Main struct {
	SdkDir string `mapstructure:"sdk_dir"`
}
