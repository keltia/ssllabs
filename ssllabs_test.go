package ssllabs

import (
	"bytes"
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

func TestClient_AnalyzeEmpty(t *testing.T) {
	Before(t)

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	an, err := c.Analyze("", false)
	require.Error(t, err)
	assert.Empty(t, an)
	assert.EqualValues(t, "empty site", err.Error())
}

// Start fresh, full restults, no options
func TestClient_AnalyzeForceFull(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts1 := map[string]string{
		"host":           site,
		"startNew":       "on",
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "on",
		"ignoreMismatch": "on",
	}

	opts2 := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "on",
		"ignoreMismatch": "on",
	}

	ftp, err := ioutil.ReadFile("testdata/ssllabs-partial.json")
	require.NoError(t, err)
	require.NotEmpty(t, ftp)

	ftc, err := ioutil.ReadFile("testdata/ssllabs-full.json")
	require.NoError(t, err)
	require.NotEmpty(t, ftc)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts1).
		Reply(200).
		BodyString(string(ftp))

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts2).
		Reply(200).
		BodyString(string(ftc))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jfta Host

	err = json.Unmarshal(ftc, &jfta)
	require.NoError(t, err)

	an, err := c.Analyze(site, true)
	require.NoError(t, err)
	assert.NotEmpty(t, an)
	assert.EqualValues(t, &jfta, an)
}

// From cache, full restults, no options
func TestClient_AnalyzeCacheFull(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	opts2 := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "on",
		"ignoreMismatch": "on",
	}

	ftc, err := ioutil.ReadFile("testdata/ssllabs-full.json")
	require.NoError(t, err)
	require.NotEmpty(t, ftc)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts2).
		Reply(200).
		BodyString(string(ftc))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jfta Host

	err = json.Unmarshal(ftc, &jfta)
	require.NoError(t, err)

	an, err := c.Analyze(site, false, opts2)
	require.NoError(t, err)
	assert.NotEmpty(t, an)
	assert.EqualValues(t, &jfta, an)
}

// From cache, full restults, with options
func TestClient_AnalyzeCacheFullOpts(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "ssllabs.com"

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "on",
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

	an, err := c.Analyze(site, false, opts)
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

	// We are removing a parameter
	mopts := map[string]string{
		"host":           site,
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "on",
		"ignoreMismatch": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(mopts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
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
	mopts := map[string]string{
		"host":           site,
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "on",
		"ignoreMismatch": "on",
	}

	fta, err := ioutil.ReadFile("testdata/ssllabs.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(mopts).
		Reply(200).
		BodyString(string(fta))

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	opts := map[string]string{"fromCache": "on"}

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

	fta, err := ioutil.ReadFile("testdata/ssllabs-endp.json")
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

	fta, err := ioutil.ReadFile("testdata/ssllabs-endp.json")
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

func TestClient_GetDetailedReport(t *testing.T) {
	Before(t)

	defer gock.Off()

	site := "www.ssllabs.com"

	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c)

	fta, err := ioutil.ReadFile("testdata/ssllabs-full.json")
	require.NoError(t, err)
	require.NotEmpty(t, fta)

	var buf bytes.Buffer

	require.NoError(t, json.Compact(&buf, fta))

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "on",
		"ignoreMismatch": "on",
		"all":            "done",
	}

	gock.New(baseURL).
		Get("/analyze").
		MatchParams(opts).
		Reply(200).
		BodyString(string(fta))

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	r, err := c.GetDetailedReport(site)
	assert.NoError(t, err)
	assert.NotEmpty(t, r)

	jr, err := json.Marshal(r)

	var buf1 bytes.Buffer

	require.NoError(t, json.Compact(&buf1, jr))

	t.Logf("%s", buf.String())
	t.Logf("%s", buf1.String())
	assert.NoError(t, err)
	//assert.EqualValues(t, buf.String(), buf1.String())
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
