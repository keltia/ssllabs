package ssllabs

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testURL = "http://localhost:1000"
)

func TestParseResults(t *testing.T) {

}

func TestAddQueryParameters(t *testing.T) {
	p := AddQueryParameters("", map[string]string{"": ""})
	assert.Equal(t, "?=", p)
}

func TestAddQueryParameters_2(t *testing.T) {
	p := AddQueryParameters("foo", map[string]string{"bar": "baz"})
	assert.Equal(t, "foo?bar=baz", p)
}

func TestPrepareRequest(t *testing.T) {
	c, err := NewClient(Config{BaseURL: testURL})
	require.NoError(t, err)

	opts := map[string]string{}
	req := c.prepareRequest("GET", "foo", opts)

	assert.NotNil(t, req)
	assert.IsType(t, (*http.Request)(nil), req)

	res, _ := url.Parse(testURL + "/foo")
	assert.Equal(t, "GET", req.Method)
	assert.EqualValues(t, res, req.URL)
}
