// subr.go

/*
Package ssllabs contains SSLLabs-related functions.
*/
package ssllabs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func myRedirect(req *http.Request, via []*http.Request) error {
	return nil
}

// AddQueryParameters adds query parameters to the URL.
func AddQueryParameters(baseURL string, queryParams map[string]string) string {
	params := url.Values{}
	if len(queryParams) == 0 {
		return baseURL
	}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// prepareRequest insert all pre-defined stuff
func (c *Client) prepareRequest(method, what string, opts map[string]string) (req *http.Request) {
	endPoint := fmt.Sprintf("%s/%s", c.baseurl, what)

	baseURL := AddQueryParameters(endPoint, opts)
	c.debug("Options:\n%v", opts)
	c.debug("baseURL: %s", baseURL)

	req, _ = http.NewRequest(method, baseURL, nil)

	c.debug("req=%#v", req)

	return
}

func (c *Client) callAPI(what, sbody string, opts map[string]string) ([]byte, error) {
	var body []byte

	retry := 0

	c.debug("callAPI")
	c.debug("clt=%#v", c.client)
	c.debug("opts=%v", opts)

	req := c.prepareRequest("GET", what, opts)
	if req == nil {
		return []byte{}, fmt.Errorf("nil req")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.debug("err=%#v", err)
		return []byte{}, errors.Wrap(err, "1st call")
	}
	defer resp.Body.Close()

	c.debug("resp=%#v", resp)

	if resp.StatusCode == http.StatusOK {

		c.debug("status OK")
		c.debug("read body")

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "body read, retry=%d", retry)
		}
		return body, errors.Wrapf(err, "status: %v body: %q", resp.Status, body)
	}
	c.debug("NOK")
	return []byte{}, errors.Wrapf(err, "status: %d", resp.StatusCode)
}

// ParseResults unmarshals the json payload
func ParseResults(content []byte) (r []Host, err error) {
	var data []Host

	err = json.Unmarshal(content, &data)
	return data, errors.Wrap(err, "unmarshal")
}

func mergeOptions(opts, o map[string]string) map[string]string {
	for i, opt := range o {
		// "" means delete
		if opt != "" {
			opts[i] = opt
		} else {
			delete(opts, i)
		}
	}
	return opts
}
