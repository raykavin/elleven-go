package elleven

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	defaultAPIPort  = "45715"
	defaultAuthPort = "45700"
	defaultTimeout  = 30 * time.Second
	sdkVersion      = "1.0.0"
	userAgent       = "elleven-go-sdk/" + sdkVersion
)

// Config holds all configuration options for the ERP Elleven SDK client.
type Config struct {
	// BaseURL is the base URL of the ERP Elleven instance, WITHOUT port.
	// Example: "https://erp.yourcompany.com" or "http://192.168.1.10"
	//
	// Used together with APIPort and AuthPort to construct APIURL and AuthURL.
	// Ignored if APIURL (or AuthURL) is set directly.
	BaseURL string

	// APIPort is the port for all non-auth API calls. Defaults to "45715".
	// Ignored if APIURL is set directly.
	APIPort string

	// AuthPort is the port for authentication endpoints. Defaults to "45700".
	// Ignored if AuthURL is set directly.
	AuthPort string

	// APIURL is the full base URL (scheme + host + port) used for API calls.
	// When set, BaseURL and APIPort are ignored for API calls.
	// Example: "http://erp.example.com:45715"
	//
	// Useful for testing (httptest.Server) or non-standard deployments.
	APIURL string

	// AuthURL is the full base URL (scheme + host + port) used for auth calls.
	// When set, BaseURL and AuthPort are ignored for auth calls.
	// Example: "http://erp.example.com:45700"
	AuthURL string

	// AccessToken is an optional pre-configured Bearer token.
	AccessToken string

	// HTTPClient allows injecting a custom *http.Client for transport-level
	// customization (TLS, proxies, etc.). If nil, a default client is created
	// using the Timeout value.
	HTTPClient *http.Client

	// Timeout for HTTP requests. Defaults to 30 seconds.
	// Only applies when HTTPClient is nil.
	Timeout time.Duration

	// MaxRetries is the number of retry attempts for transient errors
	// (HTTP 429, 502, 503, 504).
	// 0 (zero value) → use default of 3 retries.
	// -1 → disable retries entirely (only one attempt is made).
	// Any positive value → use exactly that many retries.
	MaxRetries int

	// RetryWaitMin is the minimum wait between retries. Defaults to 1 second.
	RetryWaitMin time.Duration

	// RetryWaitMax is the maximum wait between retries (exponential cap).
	// Defaults to 30 seconds.
	RetryWaitMax time.Duration
}

// setDefaults fills in default values for unset Config fields.
func (c *Config) setDefaults() {
	if c.APIPort == "" {
		c.APIPort = defaultAPIPort
	}
	if c.AuthPort == "" {
		c.AuthPort = defaultAuthPort
	}
	// Derive APIURL / AuthURL from BaseURL if not set directly.
	if c.APIURL == "" && c.BaseURL != "" {
		c.APIURL = fmt.Sprintf("%s:%s", strings.TrimRight(c.BaseURL, "/"), c.APIPort)
	}
	if c.AuthURL == "" && c.BaseURL != "" {
		c.AuthURL = fmt.Sprintf("%s:%s", strings.TrimRight(c.BaseURL, "/"), c.AuthPort)
	}
	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}
	switch {
	case c.MaxRetries == 0:
		c.MaxRetries = 3 // zero value → default
	case c.MaxRetries < 0:
		c.MaxRetries = 0 // -1 → disabled
	}
	if c.RetryWaitMin == 0 {
		c.RetryWaitMin = time.Second
	}
	if c.RetryWaitMax == 0 {
		c.RetryWaitMax = 30 * time.Second
	}
}
