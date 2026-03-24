package elleven

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateContract_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/external/billing/contracts/create") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body CreateContractRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.ClientInformation.TxID != "12345678900" {
			t.Errorf("clientTxId: want '12345678900', got %q", body.ClientInformation.TxID)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"Success": true,
			"response": map[string]any{
				"contractNumber":       "C001234",
				"activationAssignment": 5678,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.CreateContract(context.Background(), CreateContractRequest{
		ClientInformation: ContractClientInfo{TxID: "12345678900"},
		ContractInformation: ContractInfo{
			ContractType:     "PF",
			CompanyPlaceTxID: "00.000.000/0001-00",
			CollectionDay:    10,
		},
		ServicesInformation: []ContractServiceInfo{
			{Code: "PLAN_FIBER_100", Quantity: 1, Price: 99.90},
		},
	})
	if err != nil {
		t.Fatalf("CreateContract: %v", err)
	}
	if !resp.Success {
		t.Error("expected Success=true")
	}
	if resp.Response.ContractNumber != "C001234" {
		t.Errorf("contractNumber: want 'C001234', got %q", resp.Response.ContractNumber)
	}
	if resp.Response.ActivationAssignment != 5678 {
		t.Errorf("activationAssignment: want 5678, got %d", resp.Response.ActivationAssignment)
	}
}

func TestApproveContract_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "contracts/approve/42") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body ApproveContractRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"personUsers": []any{},
				"proportionalities": []map[string]any{
					{
						"description":        "Proporcionalidade gerada na competência atual",
						"amount":             49.90,
						"competence":         "2025-11-01T00:00:00",
						"serviceDescription": "Mensalidade Plano X",
						"type":               1,
						"competenceType":     1,
					},
				},
			},
			"elapsedTime": 12.54,
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.ApproveContract(context.Background(), 42, ApproveContractRequest{
		ProportionalityType:           0,
		ChangeBeginningDateOnApproval: true,
		ChangeLoyaltyDateOnApproval:   true,
		ApprovalDate:                  "2025-11-15T19:31:19.544Z",
	})
	if err != nil {
		t.Fatalf("ApproveContract: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if len(resp.Data.Proportionalities) != 1 {
		t.Fatalf("expected 1 proportionality, got %d", len(resp.Data.Proportionalities))
	}
	if resp.Data.Proportionalities[0].Amount != 49.90 {
		t.Errorf("amount: want 49.90, got %f", resp.Data.Proportionalities[0].Amount)
	}
}

func TestListDocumentsByContract_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "getfilesbycontract/C001/documentationtype/1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": []map[string]any{
				{
					"title":             "Contrato Assinado",
					"description":       "<p>Contrato de prestação de serviços</p>",
					"documentationType": "Documentos Contratuais",
					"validity":          map[string]any{"begin": nil, "final": nil},
					"file":              "https://example.com/files/contrato.pdf",
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.ListDocumentsByContract(context.Background(), "C001", 1)
	if err != nil {
		t.Fatalf("ListDocumentsByContract: %v", err)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("expected 1 document, got %d", len(resp.Response))
	}
	if resp.Response[0].Title != "Contrato Assinado" {
		t.Errorf("title: want 'Contrato Assinado', got %q", resp.Response[0].Title)
	}
}

func TestUploadAttachmentToContract_Multipart(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("expected multipart content-type, got: %s", r.Header.Get("Content-Type"))
		}
		if r.URL.Query().Get("contractNumber") != "C001234" {
			t.Errorf("contractNumber: want 'C001234', got %q", r.URL.Query().Get("contractNumber"))
		}
		if r.URL.Query().Get("documentationTypeCode") != "1.01" {
			t.Errorf("documentationTypeCode: want '1.01', got %q", r.URL.Query().Get("documentationTypeCode"))
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, err := client.UploadAttachmentToContract(
		context.Background(),
		"C001234",
		"1.01",
		"contrato.pdf",
		[]byte("%PDF test content"),
	)
	if err != nil {
		t.Fatalf("UploadAttachmentToContract: %v", err)
	}
}

func TestDownloadContractAttachment_Binary(t *testing.T) {
	fileContent := []byte("binary file content here")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id") != "5" {
			t.Errorf("id: want '5', got %q", r.URL.Query().Get("id"))
		}
		if r.URL.Query().Get("contractId") != "100" {
			t.Errorf("contractId: want '100', got %q", r.URL.Query().Get("contractId"))
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(fileContent)
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	data, ct, err := client.DownloadContractAttachment(context.Background(), 5, 100)
	if err != nil {
		t.Fatalf("DownloadContractAttachment: %v", err)
	}
	if ct != "application/octet-stream" {
		t.Errorf("content-type: want 'application/octet-stream', got %q", ct)
	}
	if string(data) != string(fileContent) {
		t.Error("file content mismatch")
	}
}

func TestCreateEventualValue_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body CreateEventualValueRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.ContractID != 12345 {
			t.Errorf("contractId: want 12345, got %d", body.ContractID)
		}
		if body.UnitAmount != 10.00 {
			t.Errorf("unitAmount: want 10.00, got %f", body.UnitAmount)
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"messages": []map[string]any{
				{"message": "1 registro criado com sucesso.", "code": 0, "type": "Success"},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.CreateEventualValue(context.Background(), CreateEventualValueRequest{
		Type:        EventualValueTypeDebit,
		ContractID:  12345,
		Description: "Desconto fidelidade",
		UnitAmount:  10.00,
		Units:       1,
		MonthYear:   "2024-12-01",
	})
	if err != nil {
		t.Fatalf("CreateEventualValue: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
}
