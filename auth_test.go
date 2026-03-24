package elleven

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// tokenHandler returns a minimal valid token JSON response.
func tokenHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"access_token":  "at-xxx",
		"expires_in":    3600,
		"refresh_token": "rt-yyy",
		"token_type":    "Bearer",
	})
}

// LegacyAuthRequest

func TestAuthenticateLegacy_DefaultParams(t *testing.T) {
	var captured url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/connect/token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if !strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			t.Errorf("unexpected content-type: %s", r.Header.Get("Content-Type"))
		}
		_ = r.ParseForm()
		captured = r.Form
		tokenHandler(w, r)
	}))
	defer srv.Close()

	client, _ := NewClient(Config{APIURL: srv.URL, AuthURL: srv.URL})

	tr, err := client.AuthenticateLegacy(context.Background(), LegacyAuthRequest{
		Username: "testuser",
		Password: "testpass",
		SynData:  "syndata123",
	})
	if err != nil {
		t.Fatalf("AuthenticateLegacy: %v", err)
	}

	// defaults applied
	assertFormValue(t, captured, "grant_type", DefaultLegacyGrantType)
	assertFormValue(t, captured, "scope", DefaultLegacyScope)
	assertFormValue(t, captured, "client_id", DefaultLegacyClientID)
	assertFormValue(t, captured, "client_secret", DefaultLegacyClientSecret)

	// user-supplied fields
	assertFormValue(t, captured, "username", "testuser")
	assertFormValue(t, captured, "password", "testpass")
	assertFormValue(t, captured, "syndata", "syndata123")

	if tr.AccessToken != "at-xxx" {
		t.Errorf("access_token: want 'at-xxx', got %q", tr.AccessToken)
	}
	if tr.ExpiresIn != 3600 {
		t.Errorf("expires_in: want 3600, got %d", tr.ExpiresIn)
	}
	if tr.RefreshToken != "rt-yyy" {
		t.Errorf("refresh_token: want 'rt-yyy', got %q", tr.RefreshToken)
	}
}

func TestAuthenticateLegacy_OverrideParams(t *testing.T) {
	var captured url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		captured = r.Form
		tokenHandler(w, r)
	}))
	defer srv.Close()

	client, _ := NewClient(Config{APIURL: srv.URL, AuthURL: srv.URL})

	_, err := client.AuthenticateLegacy(context.Background(), LegacyAuthRequest{
		GrantType:    "custom_grant",
		Scope:        "custom_scope",
		ClientID:     "custom-client",
		ClientSecret: "custom-secret",
		Username:     "u",
		Password:     "p",
		SynData:      "s",
	})
	if err != nil {
		t.Fatalf("AuthenticateLegacy: %v", err)
	}

	assertFormValue(t, captured, "grant_type", "custom_grant")
	assertFormValue(t, captured, "scope", "custom_scope")
	assertFormValue(t, captured, "client_id", "custom-client")
	assertFormValue(t, captured, "client_secret", "custom-secret")
}

// ClientCredentialsRequest

func TestAuthenticateClientCredentials_DefaultParams(t *testing.T) {
	var captured url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		captured = r.Form
		tokenHandler(w, r)
	}))
	defer srv.Close()

	client, _ := NewClient(Config{APIURL: srv.URL, AuthURL: srv.URL})

	tr, err := client.AuthenticateClientCredentials(context.Background(), ClientCredentialsRequest{
		ClientID:     "my-client-id",
		ClientSecret: "my-secret",
		SynData:      "syndata456",
	})
	if err != nil {
		t.Fatalf("AuthenticateClientCredentials: %v", err)
	}

	assertFormValue(t, captured, "grant_type", DefaultClientCredentialsGrantType)
	assertFormValue(t, captured, "scope", DefaultClientCredentialsScope)
	assertFormValue(t, captured, "client_id", "my-client-id")
	assertFormValue(t, captured, "client_secret", "my-secret")
	assertFormValue(t, captured, "syndata", "syndata456")

	if tr.AccessToken != "at-xxx" {
		t.Errorf("access_token: want 'at-xxx', got %q", tr.AccessToken)
	}
}

func TestAuthenticateClientCredentials_OverrideParams(t *testing.T) {
	var captured url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		captured = r.Form
		tokenHandler(w, r)
	}))
	defer srv.Close()

	client, _ := NewClient(Config{APIURL: srv.URL, AuthURL: srv.URL})

	_, err := client.AuthenticateClientCredentials(context.Background(), ClientCredentialsRequest{
		GrantType:    "custom_cc",
		Scope:        "syngw extra_scope",
		ClientID:     "id",
		ClientSecret: "secret",
	})
	if err != nil {
		t.Fatalf("AuthenticateClientCredentials: %v", err)
	}

	assertFormValue(t, captured, "grant_type", "custom_cc")
	assertFormValue(t, captured, "scope", "syngw extra_scope")
}

// RefreshTokenRequest

func TestRefreshToken_DefaultParams(t *testing.T) {
	var captured url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		captured = r.Form
		tokenHandler(w, r)
	}))
	defer srv.Close()

	client, _ := NewClient(Config{APIURL: srv.URL, AuthURL: srv.URL})

	tr, err := client.RefreshToken(context.Background(), RefreshTokenRequest{
		RefreshToken: "old-refresh-token",
	})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}

	assertFormValue(t, captured, "grant_type", DefaultRefreshTokenGrantType)
	assertFormValue(t, captured, "client_id", DefaultLegacyClientID)
	assertFormValue(t, captured, "client_secret", DefaultLegacyClientSecret)
	assertFormValue(t, captured, "refresh_token", "old-refresh-token")

	if tr.AccessToken != "at-xxx" {
		t.Errorf("want 'at-xxx', got %q", tr.AccessToken)
	}
}

func TestRefreshToken_OverrideParams(t *testing.T) {
	var captured url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		captured = r.Form
		tokenHandler(w, r)
	}))
	defer srv.Close()

	client, _ := NewClient(Config{APIURL: srv.URL, AuthURL: srv.URL})

	_, err := client.RefreshToken(context.Background(), RefreshTokenRequest{
		GrantType:    "custom_refresh",
		ClientID:     "other-client",
		ClientSecret: "other-secret",
		RefreshToken: "tok",
	})
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}

	assertFormValue(t, captured, "grant_type", "custom_refresh")
	assertFormValue(t, captured, "client_id", "other-client")
	assertFormValue(t, captured, "client_secret", "other-secret")
}

// Error / edge cases

func TestAuthentication_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant","error_description":"Invalid credentials"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(Config{
		APIURL:     srv.URL,
		AuthURL:    srv.URL,
		MaxRetries: 0,
	})

	_, err := client.AuthenticateLegacy(context.Background(), LegacyAuthRequest{
		Username: "bad",
		Password: "bad",
		SynData:  "bad",
	})
	if err == nil {
		t.Fatal("expected error for 400 response")
	}
}

// TestAuthURL_DifferentURLs ensures that auth requests go to AuthURL and
// API requests go to APIURL when they are configured separately.
func TestAuthURL_DifferentURLs(t *testing.T) {
	var authHit, apiHit bool

	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHit = true
		resp := map[string]any{
			"access_token": "tok", "expires_in": 3600, "token_type": "Bearer",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer authSrv.Close()

	apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiHit = true
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"data": []any{}, "totalRecords": 0, "page": 1, "pageSize": 1, "totalPages": 0,
			},
		})
	}))
	defer apiSrv.Close()

	client, err := NewClient(Config{
		APIURL:  apiSrv.URL,
		AuthURL: authSrv.URL,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tr, err := client.AuthenticateClientCredentials(context.Background(), ClientCredentialsRequest{
		ClientID:     "id",
		ClientSecret: "secret",
		SynData:      "data",
	})
	if err != nil {
		t.Fatalf("AuthenticateClientCredentials: %v", err)
	}
	client.SetTokenResponse(tr)

	if !authHit {
		t.Error("auth server was not hit")
	}

	_, err = client.ListCustomers(context.Background(), PaginationParams{Page: 1, PageSize: 1})
	if err != nil {
		t.Fatalf("ListCustomers: %v", err)
	}
	if !apiHit {
		t.Error("api server was not hit")
	}
}

// TestURL_EncodingInFormValues ensures special characters in form values are encoded.
func TestURL_EncodingInFormValues(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		tokenHandler(w, r)
	}))
	defer srv.Close()

	client, _ := NewClient(Config{
		APIURL:     srv.URL,
		AuthURL:    srv.URL,
		MaxRetries: 0,
	})
	_, _ = client.AuthenticateLegacy(context.Background(), LegacyAuthRequest{
		Username: "user@domain.com",
		Password: "p@ss!word#1",
		SynData:  "token",
	})

	values, err := url.ParseQuery(gotBody)
	if err != nil {
		t.Fatalf("parsing form body: %v", err)
	}
	if values.Get("username") != "user@domain.com" {
		t.Errorf("username not properly decoded: %q", values.Get("username"))
	}
	if values.Get("password") != "p@ss!word#1" {
		t.Errorf("password not properly decoded: %q", values.Get("password"))
	}
}

// TestTokenResponse_ExpiryCalculation checks that expiry is set correctly.
func TestTokenResponse_ExpiryCalculation(t *testing.T) {
	client, _ := NewClient(Config{BaseURL: "http://localhost"})

	before := time.Now()
	client.SetTokenResponse(&TokenResponse{
		AccessToken: "tok",
		ExpiresIn:   1800, // 30 minutes
	})
	after := time.Now()

	exp := client.TokenExpiresAt()
	minExp := before.Add(1800 * time.Second)
	maxExp := after.Add(1800 * time.Second)

	if exp.Before(minExp) || exp.After(maxExp) {
		t.Errorf("expiry %v not in range [%v, %v]", exp, minExp, maxExp)
	}
}

// Helpers

func assertFormValue(t *testing.T, form url.Values, key, want string) {
	t.Helper()
	if got := form.Get(key); got != want {
		t.Errorf("form[%q]: want %q, got %q", key, want, got)
	}
}
