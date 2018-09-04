package ssllabs

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/goware/httpmock"
	"github.com/h2non/gock"
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
	conf := Config{BaseURL: testURL, Log: 1}
	c, err := NewClient(conf)

	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotEmpty(t, c)

	assert.Equal(t, testURL, c.baseurl)
}

func TestNewClient4(t *testing.T) {
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

	fte, err := ioutil.ReadFile("testdata/emptyanalyze.json")
	require.NoError(t, err)
	require.NotEmpty(t, fte)

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

func TestClient_Analyze2(t *testing.T) {
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

func TestClient_Analyze3(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	opts := map[string]string{"foo": "bar"}

	an, err := c.Analyze("", opts)
	require.Error(t, err)
	assert.Empty(t, an)
}

func TestClient_GetStatusCodes(t *testing.T) {
	Before(t)

	ftr, err := ioutil.ReadFile("testdata/statuscodes.json")
	require.NoError(t, err)
	require.NotEmpty(t, ftr)

	gock.New(testURL).
		Get("/getStatusCodes").
		Reply(200).
		BodyString(string(ftr))

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	sc, err := c.GetStatusCodes()
	require.NoError(t, err)
	assert.NotEmpty(t, sc)
}

func TestClient_Info(t *testing.T) {
	defer gock.Off()

	fti, err := ioutil.ReadFile("testdata/info.json")
	require.NoError(t, err)
	require.NotEmpty(t, fti)

	Before(t)
	gock.New(testURL).
		Get("/info").
		Reply(200).
		BodyString(string(fti))

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	info, err := c.Info()
	require.NoError(t, err)
	assert.NotEmpty(t, info)
}

func TestClient_GetGrade(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	grade, err := c.GetGrade("lbl.gov")
	assert.Error(t, err)
	assert.Empty(t, grade)
}

func TestClient_GetGrade2(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	opts := map[string]string{"foo": "bar"}

	grade, err := c.GetGrade("lbl.gov", opts)
	assert.Error(t, err)
	assert.Empty(t, grade)
}

func TestClient_GetEndpointData(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	grade, err := c.GetEndpointData("lbl.gov")
	assert.Error(t, err)
	assert.Empty(t, grade)
}

func TestClient_GetEndpointData2(t *testing.T) {
	Before(t)
	BeforeAPI(t)

	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	opts := map[string]string{"foo": "bar"}

	grade, err := c.GetGrade("lbl.gov", opts)
	assert.Error(t, err)
	assert.Empty(t, grade)
}

func TestVersion(t *testing.T) {
	v := Version()
	assert.Equal(t, MyVersion, v)
}
