package ssllabs

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/goware/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testURL = "http://localhost:10000"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewClient2(t *testing.T) {
	conf := Config{BaseURL: testURL}
	c, err := NewClient(conf)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotEmpty(t, c)

	assert.Equal(t, testURL, c.baseurl)
}

func TestNewClient3(t *testing.T) {
	conf := Config{BaseURL: testURL, Log: 2}
	c, err := NewClient(conf)

	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotEmpty(t, c)

	assert.Equal(t, testURL, c.baseurl)
}

func Before(t *testing.T) {
	os.Unsetenv("http_proxy")
	os.Unsetenv("https_proxy")
	os.Unsetenv("all_proxy")
}

var (
	mockService *httpmock.MockHTTPServer
)

func BeforeAPI(t *testing.T) {
	var err error

	if mockService == nil {
		// new mocking server
		t.Log("starting mock...")
		mockService = httpmock.NewMockHTTPServer("localhost:10000")
	}

	require.NotNil(t, mockService)

	// define request->response pairs
	request1, _ := url.Parse(testURL + "/analyze?host=lbl.gov")
	request2, _ := url.Parse(testURL + "/getStatusCodes")

	ftr, err := ioutil.ReadFile("testdata/statuscodes.json")
	assert.NoError(t, err)

	aresp := []httpmock.MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    request1,
			},
			Response: httpmock.Response{
				StatusCode: 200,
				Body:       "done",
			},
		},
		{
			Request: http.Request{
				Method: "GET",
				URL:    request2,
			},
			Response: httpmock.Response{
				StatusCode: 200,
				Body:       string(ftr),
			},
		},
	}

	mockService.AddResponses(aresp)
	//t.Logf("respmap=%v", mockService.ResponseMap)
}

func TestClient_GetStatusCodes(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	sc, err := c.GetStatusCodes()
	require.NoError(t, err)
	assert.NotEmpty(t, sc)
}

func TestVersion(t *testing.T) {
	v := Version()
	assert.Equal(t, MyVersion, v)
}
