package elleven

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/raykavin/ellevensdk/internal/retry"
)

// Client is the main ERP Elleven SDK client.
// It is safe for concurrent use by multiple goroutines.
//
// Create a new Client with NewClient, then authenticate using one of the
// Authenticate* methods before making API calls.
type Client struct {
	config     Config
	httpClient *http.Client

	tokenMu sync.RWMutex
	token   tokenState
}

type tokenState struct {
	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

// NewClient creates and returns a new ERP Elleven SDK client.
// Returns an error if the configuration is invalid.
//
// Either BaseURL or both APIURL and AuthURL must be provided.
func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" && strings.TrimSpace(cfg.APIURL) == "" {
		return nil, fmt.Errorf("elleven: Config.BaseURL is required (or set Config.APIURL + Config.AuthURL)")
	}
	cfg.setDefaults()

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
		}
	}

	c := &Client{
		config:     cfg,
		httpClient: httpClient,
	}

	if cfg.AccessToken != "" {
		c.token.accessToken = cfg.AccessToken
	}

	return c, nil
}

// SetAccessToken sets the Bearer token used for API requests.
// This replaces any existing token.
func (c *Client) SetAccessToken(token string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token.accessToken = token
}

// SetTokenResponse stores a full token response (access + refresh + expiry).
// Call this after a successful Authenticate* call.
func (c *Client) SetTokenResponse(tr *TokenResponse) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token.accessToken = tr.AccessToken
	c.token.refreshToken = tr.RefreshToken
	if tr.ExpiresIn > 0 {
		c.token.expiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	}
}

// GetAccessToken returns the current access token.
func (c *Client) GetAccessToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token.accessToken
}

// GetRefreshToken returns the current refresh token.
func (c *Client) GetRefreshToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token.refreshToken
}

// TokenExpiresAt returns the time when the current access token expires.
// Returns zero time if no expiry is known.
func (c *Client) TokenExpiresAt() time.Time {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token.expiresAt
}

// apiURL constructs a full URL for an API endpoint path.
func (c *Client) apiURL(path string) string {
	return strings.TrimRight(c.config.APIURL, "/") + path
}

// authURL constructs a full URL for an authentication endpoint path.
func (c *Client) authURL(path string) string {
	return strings.TrimRight(c.config.AuthURL, "/") + path
}

// do executes an HTTP request, injecting the Bearer token and User-Agent,
// and applying retry logic for transient errors.
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	c.tokenMu.RLock()
	token := c.token.accessToken
	c.tokenMu.RUnlock()

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("User-Agent", userAgent)
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}

	var (
		lastErr  error
		lastResp *http.Response
	)

	// We need to buffer the body for retries.
	var bodyBytes []byte
	if req.Body != nil && req.Body != http.NoBody {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("elleven: reading request body for retry: %w", err)
		}
		req.Body.Close()
	}

	maxAttempts := c.config.MaxRetries + 1
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			wait := retry.Wait(attempt, c.config.RetryWaitMin, c.config.RetryWaitMax)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
		}

		// Restore body for each attempt.
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			req.ContentLength = int64(len(bodyBytes))
		}

		resp, err := c.httpClient.Do(req.WithContext(ctx))
		if err != nil {
			lastErr = err
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			// Retry on transport errors.
			continue
		}

		if isRetryableStatus(resp.StatusCode) && attempt < maxAttempts-1 {
			_ = resp.Body.Close()
			lastResp = resp
			lastErr = fmt.Errorf("elleven: server returned %d", resp.StatusCode)
			continue
		}

		return resp, nil
	}

	if lastResp != nil {
		return lastResp, nil
	}
	return nil, lastErr
}

// doJSON performs an HTTP request with an optional JSON body and decodes
// the response into result (if non-nil).
func (c *Client) doJSON(ctx context.Context, method, rawURL string, body, result any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("elleven: marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, reqBody)
	if err != nil {
		return fmt.Errorf("elleven: creating request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("elleven: reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("elleven: decoding response (status %d): %w (body: %.200s)",
				resp.StatusCode, err, string(respBody))
		}
	}

	return nil
}

// doForm performs an HTTP request with an application/x-www-form-urlencoded body.
func (c *Client) doForm(ctx context.Context, method, rawURL string, form url.Values, result any) error {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("elleven: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("elleven: reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("elleven: decoding response: %w (body: %.200s)", err, string(respBody))
		}
	}

	return nil
}

// doMultipart performs a multipart/form-data file upload.
// fieldName is the form field name for the file (e.g., "File").
// fileName is the name of the file.
// fileData is the file content.
// extraParams are additional form fields to include.
func (c *Client) doMultipart(
	ctx context.Context,
	rawURL string,
	fieldName, fileName string,
	fileData []byte,
	result any,
) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return fmt.Errorf("elleven: creating form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return fmt.Errorf("elleven: writing file data: %w", err)
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, &buf)
	if err != nil {
		return fmt.Errorf("elleven: creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("elleven: reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("elleven: decoding response: %w", err)
		}
	}

	return nil
}

// doDownload performs an HTTP GET and returns the raw response bytes and content type.
// Used for binary responses such as PDFs and file downloads.
func (c *Client) doDownload(ctx context.Context, rawURL string) (data []byte, contentType string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("elleven: creating request: %w", err)
	}
	req.Header.Set("Accept", "*/*")

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", parseAPIError(resp.StatusCode, body)
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("elleven: reading response body: %w", err)
	}

	return data, resp.Header.Get("Content-Type"), nil
}

// isRetryableStatus returns true for HTTP status codes that should be retried.
func isRetryableStatus(status int) bool {
	switch status {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

// buildURL constructs a URL with query parameters from a map.
func buildURL(base string, params map[string]string) string {
	if len(params) == 0 {
		return base
	}
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	q := u.Query()
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}
