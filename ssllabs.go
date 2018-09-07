package ssllabs // import "github.com/keltia/ssllabs"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keltia/proxy"
	"github.com/pkg/errors"
)

/*
SSLabs API v3

https://github.com/ssllabs/ssllabs-scan/blob/master/ssllabs-api-docs-v3.md

GET only, no POST
*/

const (
	baseURL = "https://api.ssllabs.com/api/v3"

	// DefaultWait is the timeout
	DefaultWait = 10 * time.Second

	// DefaultRetry is the number of retries we allow
	DefaultRetry = 5

	// MyVersion is the API version
	MyVersion = "0.9.0"

	// MyName is the name used for the configuration
	MyName = "ssllabs"
)

// Client is the main datatype for requests
type Client struct {
	baseurl   string
	level     int
	timeout   time.Duration
	retries   int
	proxyauth string

	client *http.Client
}

// Config is for the client configuration
type Config struct {
	BaseURL string
	Log     int
	Timeout int
	Retries int
}

// NewClient create the context for new connections
func NewClient(cnf ...Config) (*Client, error) {
	var c *Client

	// Set default
	if len(cnf) == 0 {
		c = &Client{
			baseurl: baseURL,
			timeout: DefaultWait,
			retries: DefaultRetry,
		}
	} else {
		c = &Client{
			baseurl: cnf[0].BaseURL,
			level:   cnf[0].Log,
			retries: cnf[0].Retries,
			timeout: toDuration(cnf[0].Timeout) * time.Second,
		}

		if cnf[0].Timeout == 0 {
			c.timeout = DefaultWait
		} else {
			c.timeout = time.Duration(cnf[0].Timeout) * time.Second
		}

		// Ensure proper default
		if c.retries == 0 {
			c.retries = DefaultRetry
		}
		// Ensure we have the API endpoint right
		if c.baseurl == "" {
			c.baseurl = baseURL
		}

		c.debug("got cnf: %#v", cnf[0])
	}

	c.verbose("client created")
	// We do not care whether it fails or not, if it does, just no proxyauth.
	proxyauth, _ := proxy.SetupProxyAuth()

	// Save it
	c.proxyauth = proxyauth
	c.debug("got proxyauth: %s", c.proxyauth)

	_, trsp := proxy.SetupTransport(c.baseurl)
	c.client = &http.Client{
		Transport:     trsp,
		Timeout:       c.timeout,
		CheckRedirect: myRedirect,
	}
	c.debug("newclient: c=%#v", c)

	return c, nil
}

// Info implements the Info() API call
func (c *Client) Info() (*Info, error) {
	// No parameter
	opts := map[string]string{}
	raw, err := c.callAPI("info", "", opts)
	if err != nil {
		return &Info{}, errors.Wrap(err, "Info")
	}

	var li Info

	err = json.Unmarshal(raw, &li)
	return &li, errors.Wrapf(err, "Info - %v", string(raw))
}

// GetGrade is the basic call — equal to getEndpointData and extracting just the grade.
func (c *Client) GetGrade(site string, myopts ...map[string]string) (string, error) {
	if site == "" {
		return "Z", errors.New("empty site")
	}

	opts := map[string]string{"all": ""}

	// Override default options
	if myopts != nil {
		for _, o := range myopts {
			opts = mergeOptions(opts, o)
		}
	}

	lr, err := c.Analyze(site, false, myopts...)
	if err != nil {
		return "Z", errors.Wrap(err, "GetGrade")
	}

	if len(lr.Endpoints) != 0 {
		if lr.Endpoints[0].StatusMessage != "Ready" {
			return "Z", fmt.Errorf("error: %s", lr.Endpoints[0].StatusMessage)
		}
		return lr.Endpoints[0].Grade, errors.Wrapf(err, "GetGrade - %v", lr)
	}
	return "Z", errors.New("no endpoint")
}

// GetDetailedReport returns the full report
func (c *Client) GetDetailedReport(site string) (Host, error) {
	return Host{}, nil
}

// Analyze submit the given host for checking
func (c *Client) Analyze(site string, force bool, myopts ...map[string]string) (*Host, error) {
	var (
		raw []byte
		err error
		lr  Host
	)

	// Default parameters
	opts := map[string]string{
		"host":           site,
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}

	if site == "" {
		return &Host{}, errors.New("empty site")
	}

	// Override default options
	if myopts != nil {
		for _, o := range myopts {
			opts = mergeOptions(opts, o)
		}
	}

	c.debug("opts=%v", opts)

	// Trigger the analyze
	if force {
		opts["startNew"] = "on"
		opts["all"] = "done"

		raw, err := c.callAPI("analyze", "", opts)
		if err != nil {
			return &Host{}, errors.Wrap(err, "analyze/trigger")
		}

		// Have a look at the body
		c.debug("raw=%v", string(raw))
	} else {
		opts["fromCache"] = "on"
	}

	retry := 0
	for {
		if retry >= c.retries {
			return &Host{}, fmt.Errorf("retries exceeded raw=%v", string(raw))
		}

		raw, err = c.callAPI("analyze", "", opts)
		if err != nil {
			return &Host{}, errors.Wrap(err, "analyze/loop")
		}

		err = json.Unmarshal(raw, &lr)
		if err != nil {
			return &Host{}, errors.Wrapf(err, "analyze/unmarshal: %s", string(raw))
		}

		c.debug("lr=%#v", lr)
		c.debug("raw=%v", string(raw))

		// End of analysis
		if lr.Status == "READY" || lr.Status == "ERROR " {
			c.debug("out-of-loop")
			break
		}

		c.debug("loop")
		time.Sleep(2 * time.Second)
		retry++
	}
	return &lr, errors.Wrapf(err, "analyze/end: %s", string(raw))
}

// GetEndpointData returns the endpoint data, no analyze run if not available
func (c *Client) GetEndpointData(site string, myopts ...map[string]string) (*Endpoint, error) {
	// Default parameters
	opts := map[string]string{
		"host":      site,
		"fromCache": "on",
	}

	if site == "" {
		return &Endpoint{}, errors.New("empty site")
	}

	// Override default options
	if myopts != nil {
		for _, o := range myopts {
			opts = mergeOptions(opts, o)
		}
	}

	raw, err := c.callAPI("getEndpointData", "", opts)
	if err != nil {
		return &Endpoint{}, errors.Wrap(err, "GetEndpointData")
	}

	var le Endpoint

	err = json.Unmarshal(raw, &le)
	return &le, errors.Wrapf(err, "GetEndpointData - %v", raw)
}

// GetStatusCodes returns all codes & their translation
func (c *Client) GetStatusCodes() (*StatusCodes, error) {
	// No parameters
	opts := map[string]string{}

	raw, err := c.callAPI("getStatusCodes", "", opts)
	if err != nil {
		return &StatusCodes{}, errors.Wrap(err, "GetStatusCodes")
	}

	var sc StatusCodes

	err = json.Unmarshal(raw, &sc)
	return &sc, errors.Wrapf(err, "GetStatusCodes - %v", string(raw))
}

// Version returns the API wrapper info
func Version() string {
	return MyVersion
}
