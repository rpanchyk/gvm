package clients

type HttpClient interface {
	Get(url string) (string, error)
}
