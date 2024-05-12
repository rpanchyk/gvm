package lister

import (
	"testing"

	"github.com/rpanchyk/gvm/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedHttpClient struct {
	mock.Mock
}

type mockedListCacher struct {
	mock.Mock
}

func (m *mockedHttpClient) Get(url string) (string, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *mockedListCacher) Get() ([]models.Sdk, error)   { return nil, nil }
func (m *mockedListCacher) Save(sdks []models.Sdk) error { return nil }

func TestDefaultListFetcher_Fetch(t *testing.T) {
	config := &models.Config{ReleaseURL: "http://url"}

	httpClient := new(mockedHttpClient)
	response := `<tr class=" ">
	<td class="filename"><a class="download" href="/dl/go1.22.3.windows-amd64.zip">go1.22.3.windows-amd64.zip</a></td>
	<td>Archive</td>
	<td>Windows</td>
	<td>amd64</td>
	<td></td>
	<td><tt>77299d1791a68f7da816bde7d7dfef1cbfff71e3</tt></td>
  </tr>`
	httpClient.On("Get", config.ReleaseURL).Return(response, nil)

	listCacher := new(mockedListCacher)
	listCacher.On("Get").Return(nil, nil)

	listFetcher := NewDefaultListFetcher(config, httpClient, listCacher)

	sdks, err := listFetcher.Fetch()
	if err != nil {
		t.Errorf("DefaultListFetcher fetch error: %s", err)
	}

	assert.Len(t, sdks, 1)
	assert.Contains(t, sdks, models.Sdk{
		URL:      "http://url/go1.22.3.windows-amd64.zip",
		FilePath: "go1.22.3.windows-amd64.zip",
		Version:  "1.22.3",
		Os:       "windows",
		Arch:     "amd64",
	})
}
