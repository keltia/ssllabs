package ssllabs

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
