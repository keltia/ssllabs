package ssllabs // import "github.com/keltia/ssllabs"

import (
	"encoding/json"
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
	MyVersion = "0.0.1"

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
func (c *Client) Info() (*LabsInfo, error) {
	// No parameter
	opts := map[string]string{}
	raw, err := c.callAPI("info", "", opts)
	if err != nil {
		return &LabsInfo{}, errors.Wrap(err, "Info")
	}

	var li LabsInfo

	err = json.Unmarshal(raw, &li)
	return &li, errors.Wrapf(err, "Info - %v", raw)
}

// GetGrade is the basic call — equal to getEndpointData and extracting just the grade.
func (c *Client) GetGrade(site string, myopts ...map[string]string) (string, error) {
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "on",
		"ignoreMismatch": "on",
	}

	if site == "" {
		return "", errors.New("empty site")
	}

	// Override default options
	if myopts != nil {
		for _, o := range myopts {
			opts = mergeOptions(opts, o)
		}
	}

	raw, err := c.callAPI("getEndpointData", "", opts)
	if err != nil {
		return "", errors.Wrap(err, "GetGrade")
	}

	var lr LabsReport

	err = json.Unmarshal(raw, &lr)
	if err != nil {
		return "", err
	}

	if len(lr.Endpoints) != 0 {
		return lr.Endpoints[0].Grade, errors.Wrapf(err, "GetGrade - %v", raw)
	}

	return "", nil
}

// GetDetailedReport returns the full report
func (c *Client) GetDetailedReport(site string) (LabsReport, error) {
	return LabsReport{}, nil
}

// Analyze submit the given host for checking
func (c *Client) Analyze(site string, myopts ...map[string]string) (*LabsReport, error) {
	// Default parameters
	opts := map[string]string{
		"host":           site,
		"all":            "done",
		"publish":        "off",
		"maxAge":         "24",
		"fromCache":      "off",
		"ignoreMismatch": "on",
	}

	if site == "" {
		return &LabsReport{}, errors.New("empty site")
	}

	// Override default options
	if myopts != nil {
		for _, o := range myopts {
			opts = mergeOptions(opts, o)
		}
	}

	raw, err := c.callAPI("analyze", "", opts)
	if err != nil {
		return &LabsReport{}, errors.Wrap(err, "Analyze")
	}

	var lr LabsReport

	err = json.Unmarshal(raw, &lr)
	return &lr, errors.Wrapf(err, "Analyze - %v", raw)
}

// GetEndpointData returns the endpoint data, no analyze run if not available
func (c *Client) GetEndpointData(site string, myopts ...map[string]string) (*LabsEndpoint, error) {
	// Default parameters
	opts := map[string]string{
		"host":      site,
		"fromCache": "on",
	}

	if site == "" {
		return &LabsEndpoint{}, errors.New("empty site")
	}

	// Override default options
	if myopts != nil {
		for _, o := range myopts {
			opts = mergeOptions(opts, o)
		}
	}

	raw, err := c.callAPI("getEndpointData", "", opts)
	if err != nil {
		return &LabsEndpoint{}, errors.Wrap(err, "GetEndpointData")
	}

	var le LabsEndpoint

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
