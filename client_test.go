package elleven

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Helpers

// newTestClient creates a Client pointing at the given test server URL.
// Both APIURL and AuthURL are set to the same test server so a single
// httptest.Server can handle all requests.
func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	client, err := NewClient(Config{
		APIURL:       serverURL,
		AuthURL:      serverURL,
		MaxRetries:   -1, // disable retries in unit tests for speed
		Timeout:      5 * time.Second,
		RetryWaitMin: time.Millisecond,
		RetryWaitMax: 5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetAccessToken("test-token")
	return client
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// NewClient

func TestNewClient_MissingBaseURL(t *testing.T) {
	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error for missing BaseURL, got nil")
	}
	if !strings.Contains(err.Error(), "BaseURL") {
		t.Errorf("error should mention BaseURL, got: %v", err)
	}
}

func TestNewClient_Defaults(t *testing.T) {
	client, err := NewClient(Config{BaseURL: "http://localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_WithAPIURL(t *testing.T) {
	client, err := NewClient(Config{
		APIURL:  "http://localhost:45715",
		AuthURL: "http://localhost:45700",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_WithAccessToken(t *testing.T) {
	client, err := NewClient(Config{
		BaseURL:     "http://localhost",
		AccessToken: "my-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := client.GetAccessToken(); got != "my-token" {
		t.Errorf("expected token 'my-token', got %q", got)
	}
}

// Token management

func TestSetTokenResponse(t *testing.T) {
	client, _ := NewClient(Config{BaseURL: "http://localhost"})

	tr := &TokenResponse{
		AccessToken:  "access-abc",
		RefreshToken: "refresh-xyz",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}
	client.SetTokenResponse(tr)

	if got := client.GetAccessToken(); got != "access-abc" {
		t.Errorf("access token: expected 'access-abc', got %q", got)
	}
	if got := client.GetRefreshToken(); got != "refresh-xyz" {
		t.Errorf("refresh token: expected 'refresh-xyz', got %q", got)
	}

	expiresAt := client.TokenExpiresAt()
	if expiresAt.IsZero() {
		t.Error("expected non-zero expiry time")
	}
	// Should expire ~1 hour from now
	diff := time.Until(expiresAt)
	if diff < 50*time.Minute || diff > 70*time.Minute {
		t.Errorf("expiry time out of range: %v", diff)
	}
}

// Authorization header injection

func TestAuthorizationHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true, "messages": nil, "response": nil,
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	client.SetAccessToken("Bearer-token-123")

	ctx := context.Background()
	// ListCustomers will make a GET request, we just need any request
	_, _ = client.ListCustomers(ctx, PaginationParams{Page: 1, PageSize: 1})

	want := "Bearer Bearer-token-123"
	if gotAuth != want {
		t.Errorf("Authorization header: want %q, got %q", want, gotAuth)
	}
}

// Error parsing

func TestAPIError_4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"messages": []map[string]any{
				{"message": "Token inválido", "code": 401, "type": "Error"},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, err := client.ListCustomers(context.Background(), PaginationParams{Page: 1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
	if len(apiErr.Messages) == 0 {
		t.Error("expected at least one message")
	}
	if apiErr.Messages[0].Message != "Token inválido" {
		t.Errorf("unexpected message: %q", apiErr.Messages[0].Message)
	}

	if !IsUnauthorized(err) {
		t.Error("IsUnauthorized should return true")
	}
}

func TestAPIError_5xx(t *testing.T) {
	// Return 503 twice then success (retry logic disabled in test client)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, err := client.ListCustomers(context.Background(), PaginationParams{Page: 1})
	if err == nil {
		t.Fatal("expected error for 503")
	}
	if !IsServerError(err) {
		t.Errorf("IsServerError should return true for 503, got false. err: %v", err)
	}
}

func TestIsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, err := client.GetPersonByTxID(context.Background(), "12345678900")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound should return true for 404")
	}
}

// Context cancellation

func TestContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		writeJSON(w, http.StatusOK, map[string]any{"success": true})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.ListCustomers(ctx, PaginationParams{Page: 1})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

// User-Agent header

func TestUserAgentHeader(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true, "messages": nil, "response": map[string]any{
				"data": []any{}, "totalRecords": 0, "page": 1, "pageSize": 1, "totalPages": 0,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, _ = client.ListCustomers(context.Background(), PaginationParams{Page: 1, PageSize: 1})

	if !strings.HasPrefix(gotUA, "elleven-go-sdk/") {
		t.Errorf("User-Agent should start with 'elleven-go-sdk/', got: %q", gotUA)
	}
}
