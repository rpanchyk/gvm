package defaulter

type Defaulter interface {
	Default(version string) error
}
