package ssllabs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

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
	conf := Config{}
	c, err := NewClient(conf)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotEmpty(t, c)

	assert.Equal(t, baseURL, c.baseurl)
}

func TestNewClient3(t *testing.T) {
	conf := Config{Log: 1}
	c, err := NewClient(conf)

	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotEmpty(t, c)

	assert.Equal(t, baseURL, c.baseurl)
}

func TestNewClient4(t *testing.T) {
	conf := Config{Log: 2}
	c, err := NewClient(conf)

	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotEmpty(t, c)

	assert.Equal(t, baseURL, c.baseurl)
}

func Before(t *testing.T) {
	os.Unsetenv("http_proxy")
	os.Unsetenv("https_proxy")
	os.Unsetenv("all_proxy")
}

func TestClient_Analyze(t *testing.T) {
	Before(t)

	defer gock.Off()

	// Default parameters
	opts := map[string]string{
		"host":           "",
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}
	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts).
		Reply(200)

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	an, err := c.Analyze("")
	require.Error(t, err)
	assert.Empty(t, an)
	assert.EqualValues(t, "empty site", err.Error())
}

func TestClient_Analyze2(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs-full.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jfta Host

	err = json.Unmarshal(fta, &jfta)
	require.NoError(t, err)

	an, err := c.Analyze(site)
	require.NoError(t, err)
	assert.NotEmpty(t, an)
	assert.EqualValues(t, &jfta, an)
}

func TestClient_Analyze3(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs-full.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jfta Host

	err = json.Unmarshal(fta, &jfta)
	require.NoError(t, err)

	opts["fromCache"] = "off"

	an, err := c.Analyze(site, opts)
	require.NoError(t, err)
	assert.NotEmpty(t, an)
	assert.EqualValues(t, &jfta, an)
}

func TestClient_GetStatusCodes(t *testing.T) {
	Before(t)

	ftr, err := ioutil.ReadFile("testdata/statuscodes.json")
	require.NoError(t, err)
	require.NotEmpty(t, ftr)

	gock.New(baseURL).
		Get("/getStatusCodes").
		Reply(200).
		BodyString(string(ftr))

	c, err := NewClient()
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
	Before(t)

	defer gock.Off()

	fti, err := ioutil.ReadFile("testdata/info.json")
	require.NoError(t, err)
	require.NotEmpty(t, fti)

	gock.New(baseURL).
		Get("/info").
		Reply(200).
		BodyString(string(fti))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	info, err := c.Info()
	require.NoError(t, err)
	assert.NotEmpty(t, info)
}

func TestClient_GetGradeEmpty(t *testing.T) {
	Before(t)

	defer gock.Off()

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	grade, err := c.GetGrade("")
	assert.Error(t, err)
	assert.Equal(t, "Z", grade)
}

func TestClient_GetGradeLbl(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "lbl.gov"

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}

	fta, err := ioutil.ReadFile("testdata/lbl.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	grade, err := c.GetGrade("lbl.gov")
	assert.Error(t, err)
	assert.Equal(t, "Z", grade)
}

func TestClient_GetGradeSSLLabs(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	c.level = 2
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	grade, err := c.GetGrade(site)
	require.NoError(t, err)
	assert.NotEmpty(t, grade)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetGradeSSLLabsFull(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs-full.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	c.level = 2
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	grade, err := c.GetGrade(site)
	require.NoError(t, err)
	assert.NotEmpty(t, grade)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetGradeSSLLabsOpts(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	opts["fromCache"] = "on"

	grade, err := c.GetGrade(site, opts)
	require.NoError(t, err)
	assert.NotEmpty(t, grade)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetEndpointData(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts := map[string]string{
		"host":      site,
		"fromCache": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/getEndpointData").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jfta Endpoint

	err = json.Unmarshal(fta, &jfta)
	require.NoError(t, err)

	data, err := c.GetEndpointData(site)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	assert.EqualValues(t, &jfta, data)
}

func TestClient_GetEndpointData2(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts := map[string]string{
		"host":      site,
		"fromCache": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/getEndpointData").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jfta Endpoint

	err = json.Unmarshal(fta, &jfta)
	require.NoError(t, err)

	opts["fromCache"] = "on"

	data, err := c.GetEndpointData(site, opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	assert.EqualValues(t, &jfta, data)
}

func TestClient_GetEndpointData3(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := ""

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	data, err := c.GetEndpointData(site)
	assert.Error(t, err)
	assert.Empty(t, data)
	assert.Equal(t, "empty site", err.Error())
}

func TestVersion(t *testing.T) {
	v := Version()
	assert.Equal(t, MyVersion, v)
}

func TestClient_GetDetailedReport(t *testing.T) {
	site := ""

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	r, err := c.GetDetailedReport(site)
	assert.NoError(t, err)
	assert.Empty(t, r)
}
