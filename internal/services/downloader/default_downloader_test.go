package downloader

import (
	"testing"

	"github.com/rpanchyk/gvm/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mocked struct {
	mock.Mock
}

func (m *mocked) Fetch() ([]models.Sdk, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Sdk), nil
}

func (m *mocked) Save(url, file string) error {
	args := m.Called(url, file)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func TestDefaultDownloader_Download(t *testing.T) {
	m := new(mocked)
	m.On("Fetch").Return([]models.Sdk{{URL: "http://url/go1.windows-amd64.zip", Version: "1"}})
	m.On("Save", "http://url/go1.windows-amd64.zip", "go1.windows-amd64.zip").Return(nil)

	d := NewDefaultDownloader(&models.Config{}, m, m)
	sdk, err := d.Download("1")
	if err != nil {
		t.Errorf("DefaultListFetcher fetch error: %s", err)
	}

	assert.NotNil(t, sdk)
	assert.Equal(t, *sdk, models.Sdk{
		URL:          "http://url/go1.windows-amd64.zip",
		Version:      "1",
		FilePath:     "go1.windows-amd64.zip",
		IsDownloaded: true,
	})
}
