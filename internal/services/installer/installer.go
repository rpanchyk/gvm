package installer

type Installer interface {
	Install(version string) error
}
