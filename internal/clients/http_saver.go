package clients

type HttpSaver interface {
	Save(url, file string) error
}
