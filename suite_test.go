package elleven

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// RegisterPerson

func TestRegisterPerson_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/external/integrations/thirdparty/people") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("expected JSON content-type, got %s", ct)
		}

		// Decode and verify request body
		var body RegisterPersonRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if body.Name != "João da Silva" {
			t.Errorf("name: want 'João da Silva', got %q", body.Name)
		}
		if body.TxID != "12345678900" {
			t.Errorf("txId: want '12345678900', got %q", body.TxID)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"messages": []map[string]any{
				{"message": "Registro criado com sucesso.", "code": nil, "type": "Success"},
			},
			"response": nil,
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.RegisterPerson(context.Background(), RegisterPersonRequest{
		TypeTxID:   "CPF",
		TxID:       "12345678900",
		Name:       "João da Silva",
		Email:      "joao@example.com",
		Client:     true,
		Situation:  1,
		PostalCode: "97000000",
		Street:     "Rua das Flores",
		Number:     "123",
		City:       "Santa Maria",
		State:      "RS",
	})
	if err != nil {
		t.Fatalf("RegisterPerson: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
}

// UpdatePerson

func TestUpdatePerson_Success(t *testing.T) {
	const personID = 42

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/42") {
			t.Errorf("path should contain '/42', got: %s", r.URL.Path)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success":  true,
			"messages": nil,
			"response": map[string]any{
				"peopleAddressId": 0,
				"message":         "Informações da pessoa atualizadas com sucesso!",
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.UpdatePerson(context.Background(), personID, UpdatePersonRequest{
		ID: personID,
		Email: UpdatePersonEmailRequest{
			Email:    "novo@email.com",
			EmailNfe: "nfe@email.com",
		},
	})
	if err != nil {
		t.Fatalf("UpdatePerson: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Response.Message != "Informações da pessoa atualizadas com sucesso!" {
		t.Errorf("unexpected message: %q", resp.Response.Message)
	}
}

// ListCustomers

func TestListCustomers_Pagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("page") != "2" {
			t.Errorf("page: want '2', got %q", q.Get("page"))
		}
		if q.Get("pageSize") != "10" {
			t.Errorf("pageSize: want '10', got %q", q.Get("pageSize"))
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"data": []map[string]any{
					{"id": "1", "name": "Cliente Teste", "txId": "12345678900"},
				},
				"totalRecords": 100,
				"page":         2,
				"pageSize":     10,
				"totalPages":   10,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.ListCustomers(context.Background(), PaginationParams{
		Page: 2, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("ListCustomers: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if len(resp.Response.Data) != 1 {
		t.Errorf("expected 1 customer, got %d", len(resp.Response.Data))
	}
	if resp.Response.TotalRecords != 100 {
		t.Errorf("totalRecords: want 100, got %d", resp.Response.TotalRecords)
	}
	if resp.Response.TotalPages != 10 {
		t.Errorf("totalPages: want 10, got %d", resp.Response.TotalPages)
	}
}

// GetPersonByTxID

func TestGetPersonByTxID_Success(t *testing.T) {
	const txID = "12345678900"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, txID) {
			t.Errorf("path should contain txID %q, got: %s", txID, r.URL.Path)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"id":     123456,
				"name":   "João Silva",
				"txId":   txID,
				"email":  "joao@example.com",
				"status": 1,
				"titles": []any{},
				"mainAddress": map[string]any{
					"street": "Rua das Flores",
					"city":   "Santa Maria",
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.GetPersonByTxID(context.Background(), txID)
	if err != nil {
		t.Fatalf("GetPersonByTxID: %v", err)
	}
	if resp.Response.ID != 123456 {
		t.Errorf("id: want 123456, got %d", resp.Response.ID)
	}
	if resp.Response.TxID != txID {
		t.Errorf("txId: want %q, got %q", txID, resp.Response.TxID)
	}
	if resp.Response.MainAddress.City != "Santa Maria" {
		t.Errorf("city: want 'Santa Maria', got %q", resp.Response.MainAddress.City)
	}
}

// IsHoliday

func TestIsHoliday_Holiday(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body HolidayCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.Date != "2024-05-13" {
			t.Errorf("date: want '2024-05-13', got %q", body.Date)
		}
		if body.UF != "RS" {
			t.Errorf("UF: want 'RS', got %q", body.UF)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"isHoliday": true,
				"name":      "Feriado em Santa Maria",
				"date":      "2024-05-13",
				"city":      "Santa Maria",
				"scope":     "Municipal",
				"uf":        "RS",
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.IsHoliday(context.Background(), HolidayCheckRequest{
		Date: "2024-05-13",
		UF:   "RS",
		City: "Santa Maria",
	})
	if err != nil {
		t.Fatalf("IsHoliday: %v", err)
	}
	if !resp.Response.IsHoliday {
		t.Error("expected isHoliday=true")
	}
	if resp.Response.Name != "Feriado em Santa Maria" {
		t.Errorf("unexpected holiday name: %q", resp.Response.Name)
	}
}

func TestIsHoliday_NotHoliday(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"isHoliday": false,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.IsHoliday(context.Background(), HolidayCheckRequest{
		Date: "2024-03-25",
		UF:   "RS",
		City: "Santa Maria",
	})
	if err != nil {
		t.Fatalf("IsHoliday: %v", err)
	}
	if resp.Response.IsHoliday {
		t.Error("expected isHoliday=false")
	}
}

// ListLocations

func TestListLocations_QueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("Filter") != "matriz" {
			t.Errorf("Filter: want 'matriz', got %q", q.Get("Filter"))
		}
		if q.Get("OrderBy") != "description" {
			t.Errorf("OrderBy: want 'description', got %q", q.Get("OrderBy"))
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"data": []map[string]any{
					{"id": 1, "description": "Matriz", "code": "01"},
				},
				"totalRecords": 1,
				"page":         1,
				"pageSize":     10,
				"totalPages":   1,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.ListLocations(context.Background(), PaginationParams{
		Page:     1,
		PageSize: 10,
		Filter:   "matriz",
		OrderBy:  "description",
	})
	if err != nil {
		t.Fatalf("ListLocations: %v", err)
	}
	if len(resp.Response.Data) != 1 {
		t.Errorf("expected 1 location, got %d", len(resp.Response.Data))
	}
	if resp.Response.Data[0].Description != "Matriz" {
		t.Errorf("description: want 'Matriz', got %q", resp.Response.Data[0].Description)
	}
}
