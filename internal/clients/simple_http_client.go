package clients

import (
	"fmt"
	"io"
	"net/http"
)

type SimpleHttpClient struct {
}

func (c SimpleHttpClient) Get(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error sending http request: %w", err)
	}
	defer response.Body.Close()
	// fmt.Printf("status code: %d\n", response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading http response: %w", err)
	}

	// fmt.Printf("Body: %s\n", body)
	return string(body), nil
}
