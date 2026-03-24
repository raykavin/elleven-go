package elleven

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestRetry_RecoversAfterTransientErrors verifies that the client retries on 503
// and eventually succeeds.
func TestRetry_RecoversAfterTransientErrors(t *testing.T) {
	var callCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n < 3 {
			// First two attempts fail with 503
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Third attempt succeeds
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"data":         []any{},
				"totalRecords": 0, "page": 1, "pageSize": 10, "totalPages": 0,
			},
		})
	}))
	defer srv.Close()

	client, err := NewClient(Config{
		APIURL:       srv.URL,
		AuthURL:      srv.URL,
		MaxRetries:   3,
		RetryWaitMin: 10 * time.Millisecond,
		RetryWaitMax: 50 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetAccessToken("tok")

	_, err = client.ListCustomers(context.Background(), PaginationParams{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}

	if atomic.LoadInt32(&callCount) < 3 {
		t.Errorf("expected at least 3 calls (2 failures + 1 success), got %d", callCount)
	}
}

// TestRetry_ExhaustsRetries verifies that after MaxRetries attempts, an error is returned.
func TestRetry_ExhaustsRetries(t *testing.T) {
	var callCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client, err := NewClient(Config{
		APIURL:       srv.URL,
		AuthURL:      srv.URL,
		MaxRetries:   2,
		RetryWaitMin: 5 * time.Millisecond,
		RetryWaitMax: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetAccessToken("tok")

	_, err = client.ListCustomers(context.Background(), PaginationParams{Page: 1})
	if err == nil {
		t.Fatal("expected error after exhausting retries, got nil")
	}

	// Should have attempted 1 initial + 2 retries = 3 total
	if got := atomic.LoadInt32(&callCount); got != 3 {
		t.Errorf("expected 3 attempts, got %d", got)
	}
}

// TestRetry_NoRetryOn4xx verifies that 4xx errors are NOT retried.
func TestRetry_NoRetryOn4xx(t *testing.T) {
	var callCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusBadRequest) // 400 should not retry
	}))
	defer srv.Close()

	client, err := NewClient(Config{
		APIURL:       srv.URL,
		AuthURL:      srv.URL,
		MaxRetries:   3,
		RetryWaitMin: 5 * time.Millisecond,
		RetryWaitMax: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetAccessToken("tok")

	_, _ = client.ListCustomers(context.Background(), PaginationParams{Page: 1})

	if got := atomic.LoadInt32(&callCount); got != 1 {
		t.Errorf("expected exactly 1 call for 400 error, got %d (retries should not occur)", got)
	}
}

// TestRetry_RetryOn429 verifies that 429 Too Many Requests is retried.
func TestRetry_RetryOn429(t *testing.T) {
	var callCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"data": []any{}, "totalRecords": 0, "page": 1, "pageSize": 1, "totalPages": 0,
			},
		})
	}))
	defer srv.Close()

	client, err := NewClient(Config{
		APIURL:       srv.URL,
		AuthURL:      srv.URL,
		MaxRetries:   2,
		RetryWaitMin: 5 * time.Millisecond,
		RetryWaitMax: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetAccessToken("tok")

	_, err = client.ListCustomers(context.Background(), PaginationParams{Page: 1, PageSize: 1})
	if err != nil {
		t.Fatalf("expected success after retry on 429, got: %v", err)
	}

	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

// TestRetry_DisabledWithNegativeOne verifies that MaxRetries=-1 disables retries entirely.
func TestRetry_DisabledWithNegativeOne(t *testing.T) {
	var callCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client, err := NewClient(Config{
		APIURL:     srv.URL,
		AuthURL:    srv.URL,
		MaxRetries: -1, // disable retries
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetAccessToken("tok")

	_, _ = client.ListCustomers(context.Background(), PaginationParams{Page: 1})

	// With retries disabled, exactly 1 attempt should be made
	if got := atomic.LoadInt32(&callCount); got != 1 {
		t.Errorf("expected exactly 1 call with MaxRetries=-1, got %d", got)
	}
}
