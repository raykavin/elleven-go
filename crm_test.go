package elleven

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStartSale_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "crm/startsale") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body StartSaleRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.TxID != "12345678900" {
			t.Errorf("txId: want '12345678900', got %q", body.TxID)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"protocolId": 999888,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.StartSale(context.Background(), StartSaleRequest{
		TxID:             "12345678900",
		CompanyPlaceTxID: "00.000.000/0001-00",
		ContractType:     "PF",
		ServiceProducts: []SaleServiceProduct{
			{Code: "PLAN001", Quantity: 1, Amount: 99.90},
		},
	})
	if err != nil {
		t.Fatalf("StartSale: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Response.ProtocolID != 999888 {
		t.Errorf("protocolId: want 999888, got %d", resp.Response.ProtocolID)
	}
}

func TestCancelSale_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body CancelSaleRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.ProtocolID != 999888 {
			t.Errorf("protocolId: want 999888, got %d", body.ProtocolID)
		}
		if body.Description != "Client requested cancellation" {
			t.Errorf("description: unexpected value %q", body.Description)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success":  true,
			"response": "Venda Cancelada com sucesso!",
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.CancelSale(context.Background(), CancelSaleRequest{
		ProtocolID:  999888,
		Description: "Client requested cancellation",
	})
	if err != nil {
		t.Fatalf("CancelSale: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
}

func TestAddContractServices_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body AddContractServicesRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.ContractNumber != "12345" {
			t.Errorf("contractNumber: want '12345', got %q", body.ContractNumber)
		}
		if len(body.NewServices) != 1 {
			t.Errorf("expected 1 service, got %d", len(body.NewServices))
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"Sucess":  true,
			"Message": "Inclusão de serviço realizada com sucesso.",
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	result, err := client.AddContractServices(context.Background(), AddContractServicesRequest{
		ContractNumber: "12345",
		NewServices: []NewContractService{
			{Code: "SVC001", Quantity: 1, Price: 49.90},
		},
	})
	if err != nil {
		t.Fatalf("AddContractServices: %v", err)
	}
	if !result.Success {
		t.Error("expected Success=true")
	}
}

func TestRemoveContractServices_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body RemoveContractServicesRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.ContractNumber != "12345" {
			t.Errorf("contractNumber: want '12345', got %q", body.ContractNumber)
		}
		if len(body.RemoveService) != 2 {
			t.Errorf("expected 2 items to remove, got %d", len(body.RemoveService))
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"Sucess":  true,
			"Message": "Remoção de serviço realizada com sucesso.",
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	result, err := client.RemoveContractServices(context.Background(), RemoveContractServicesRequest{
		ContractNumber: "12345",
		RemoveService: []RemoveContractItem{
			{ContractItem: 101},
			{ContractItem: 102},
		},
	})
	if err != nil {
		t.Fatalf("RemoveContractServices: %v", err)
	}
	if !result.Success {
		t.Error("expected Success=true")
	}
}

func TestVerifyViability_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body VerifyViabilityRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.FullAddress.PostalCode != "97010001" {
			t.Errorf("postalCode: want '97010001', got %q", body.FullAddress.PostalCode)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"viability": true,
				"ctos":      37,
				"ports":     11,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.VerifyViability(context.Background(), VerifyViabilityRequest{
		FullAddress: ViabilityAddress{
			Address:      "Rua das Flores",
			Number:       "123",
			Neighborhood: "Centro",
			City:         "Santa Maria",
			State:        "RS",
			PostalCode:   "97010001",
		},
		Distance: "100",
		Lead: ViabilityPerson{
			ID:   "1",
			Name: "João",
			TxID: "12345678900",
		},
	})
	if err != nil {
		t.Fatalf("VerifyViability: %v", err)
	}
	if !resp.Response.Viability {
		t.Error("expected viability=true")
	}
	if resp.Response.CTOs != 37 {
		t.Errorf("ctos: want 37, got %d", resp.Response.CTOs)
	}
}

func TestGetContractTypesAndServices_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": []map[string]any{
				{
					"code":  "1",
					"title": "PF - Composto",
					"contractTypesServiceProduct": []map[string]any{
						{"code": "1.1", "title": "Serviço de Configuração"},
					},
					"collectionDays": []int{5, 10, 15},
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.GetContractTypesAndServices(context.Background())
	if err != nil {
		t.Fatalf("GetContractTypesAndServices: %v", err)
	}
	if len(resp.Response) != 1 {
		t.Errorf("expected 1 contract type, got %d", len(resp.Response))
	}
	if resp.Response[0].Code != "1" {
		t.Errorf("code: want '1', got %q", resp.Response[0].Code)
	}
	if len(resp.Response[0].ContractTypesServiceProduct) != 1 {
		t.Errorf("expected 1 service product, got %d", len(resp.Response[0].ContractTypesServiceProduct))
	}
}
