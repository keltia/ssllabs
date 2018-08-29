package ssllabs // import "keltia/net/ssllabs"

import (
	"net/http"
	"time"

	"github.com/keltia/proxy"
)

/*
SSLabs API v3
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

type Client struct {
	baseurl   string
	level     int
	timeout   time.Duration
	retries   int
	proxyauth string

	client *http.Client
}

type Config struct {
	BaseURL string
	Log     int
	Timeout int
	Retries int
}

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

	return &Client{}, nil
}

// Info implements the Info() API call
func (c *Client) Info() (*Info, error) {
	return nil, nil
}

// GetGrade is the basic call
func (c *Client) GetGrade(site string) (string, error) {
	return "", nil
}

// GetDetailedReport returns the full report
func (c *Client) GetDetailedReport(site string) (LabsReport, error) {
	return LabsReport{}, nil
}

// Version returns the API wrapper info
func Version() string {
	return MyVersion
}
