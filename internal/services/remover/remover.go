package remover

type Remover interface {
	Remove(version string, removeDownloaded, removeInstalled bool)
}
