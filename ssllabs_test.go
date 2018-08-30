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
	request1, _ := url.Parse(testURL + "/analyze?host=")
	request2, _ := url.Parse(testURL + "/getStatusCodes")
	request3, _ := url.Parse(testURL + "/info")

	fte, err := ioutil.ReadFile("testdata/emptyanalyze.json")
	require.NoError(t, err)
	require.NotEmpty(t, fte)

	ftr, err := ioutil.ReadFile("testdata/statuscodes.json")
	require.NoError(t, err)
	require.NotEmpty(t, ftr)

	fti, err := ioutil.ReadFile("testdata/info.json")
	require.NoError(t, err)
	require.NotEmpty(t, fti)

	aresp := []httpmock.MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    request1,
			},
			Response: httpmock.Response{
				StatusCode: 200,
				Body:       string(fte),
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
		{
			Request: http.Request{
				Method: "GET",
				URL:    request3,
			},
			Response: httpmock.Response{
				StatusCode: 200,
				Body:       string(fti),
			},
		},
	}

	mockService.AddResponses(aresp)
	//t.Logf("respmap=%v", mockService.ResponseMap)
}

func TestClient_Analyze(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	an, err := c.Analyze("")
	require.Error(t, err)
	assert.Empty(t, an)
}

func TestClient_Analyze_2(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: "http://localhost:10001"})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	an, err := c.Analyze("")
	require.Error(t, err)
	assert.Empty(t, an)
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

func TestClient_Info(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	info, err := c.Info()
	require.NoError(t, err)
	assert.NotEmpty(t, info)
}

func TestVersion(t *testing.T) {
	v := Version()
	assert.Equal(t, MyVersion, v)
}
