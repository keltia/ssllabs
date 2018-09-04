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
	assert.IsType(t, ([]Host)(nil), data)
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

func TestMergeOptions(t *testing.T) {
	o1 := make(map[string]string)
	o2 := make(map[string]string)

	o3 := mergeOptions(o1, o2)
	require.Empty(t, o3)
	assert.EqualValues(t, o3, o1)
	assert.EqualValues(t, o3, o2)
}

func TestMergeOptions2(t *testing.T) {
	o1 := map[string]string{"foo": "bar"}
	o2 := make(map[string]string)

	o3 := mergeOptions(o1, o2)
	require.NotEmpty(t, o3)
	assert.EqualValues(t, o1, o3)
}

func TestMergeOptions3(t *testing.T) {
	o1 := map[string]string{"foo": "bar"}
	o2 := map[string]string{"baz": "xyzt"}
	ot := map[string]string{"baz": "xyzt", "foo": "bar"}

	o3 := mergeOptions(o1, o2)
	require.NotEmpty(t, o3)
	assert.EqualValues(t, ot, o3)
}
