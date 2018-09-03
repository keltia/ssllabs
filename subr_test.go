package ssllabs

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyRedirect(t *testing.T) {
	err := myRedirect(nil, nil)
	require.NoError(t, err)
}

func TestParseResults(t *testing.T) {
	data, err := ParseResults([]byte{})

	assert.Error(t, err)
	assert.Empty(t, data)
	assert.IsType(t, ([]LabsReport)(nil), data)
}

func TestAddQueryParameters(t *testing.T) {
	p := AddQueryParameters("", map[string]string{})
	assert.Equal(t, "", p)
}

func TestAddQueryParameters_1(t *testing.T) {
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

func TestPrepareRequest_2(t *testing.T) {
	c, err := NewClient()
	require.NoError(t, err)

	opts := map[string]string{}
	req := c.prepareRequest("GET", "foo", opts)

	assert.NotNil(t, req)
	assert.IsType(t, (*http.Request)(nil), req)

	res, _ := url.Parse(baseURL + "/foo")
	assert.Equal(t, "GET", req.Method)
	assert.EqualValues(t, res, req.URL)
}

func TestLabsErrorResponse_Error(t *testing.T) {
	var empty = "{\"errors\":null}"

	e := LabsErrorResponse{}
	msg := e.Error()

	assert.NotEmpty(t, msg)
	assert.Equal(t, empty, msg)
}

func TestLabsErrorResponse_Error2(t *testing.T) {
	var empty = "{\"errors\":[{\"Field\":\"\",\"Message\":\"\"}]}"

	e := LabsErrorResponse{ResponseErrors: []LabsError{{"", ""}}}
	msg := e.Error()

	assert.NotEmpty(t, msg)
	assert.Equal(t, empty, msg)
}

func TestLabsErrorResponse_Error3(t *testing.T) {
	var err = "{\"errors\":[{\"Field\":\"\\ufffd\\ufffd\",\"Message\":\"\"}]}"

	e := LabsErrorResponse{ResponseErrors: []LabsError{
		{string([]byte{155, 134}), ""},
	},
	}
	msg := e.Error()

	assert.NotEmpty(t, msg)
	assert.Equal(t, err, msg)
}
