package elleven

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetOpenInvoicesByTxID_Success(t *testing.T) {
	const txID = "12345678900"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "getopentitlesbytxid/"+txID) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": []map[string]any{
				{
					"id": 123456,
					"billet": map[string]any{
						"title":          "FAT000000001",
						"issueDate":      "2024-01-01",
						"expirationDate": "2024-01-31",
						"typefulLine":    "00000.00000 00000.000000 00000.00000 0 000000000000",
						"pixQRCode":      "pix-qr-code-data",
						"amount": map[string]any{
							"value":      99.90,
							"finalValue": 99.90,
							"discount":   0,
							"fine":       0,
							"interest":   0,
						},
					},
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.GetOpenInvoicesByTxID(context.Background(), txID)
	if err != nil {
		t.Fatalf("GetOpenInvoicesByTxID: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if len(resp.Response) != 1 {
		t.Fatalf("expected 1 invoice, got %d", len(resp.Response))
	}
	inv := resp.Response[0]
	if inv.ID != 123456 {
		t.Errorf("id: want 123456, got %d", inv.ID)
	}
	if inv.Billet.Title != "FAT000000001" {
		t.Errorf("title: want 'FAT000000001', got %q", inv.Billet.Title)
	}
	if inv.Billet.Amount.Value != 99.90 {
		t.Errorf("amount.value: want 99.90, got %f", inv.Billet.Amount.Value)
	}
}

func TestGetAllInvoicesByTxID_WithStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "gettitlesbytxid/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": []map[string]any{
				{
					"id":     1,
					"status": "Paga",
					"billet": map[string]any{
						"title":          "FAT000000001",
						"issueDate":      "2024-01-01",
						"expirationDate": "2024-01-31",
						"amount": map[string]any{
							"value": 99.90, "finalValue": 99.90,
						},
					},
				},
				{
					"id":     2,
					"status": "Em aberto",
					"billet": map[string]any{
						"title":          "FAT000000002",
						"issueDate":      "2024-02-01",
						"expirationDate": "2024-02-29",
						"amount": map[string]any{
							"value": 99.90, "finalValue": 99.90,
						},
					},
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.GetAllInvoicesByTxID(context.Background(), "12345678900")
	if err != nil {
		t.Fatalf("GetAllInvoicesByTxID: %v", err)
	}
	if len(resp.Response) != 2 {
		t.Fatalf("expected 2 invoices, got %d", len(resp.Response))
	}
	if resp.Response[0].Status != "Paga" {
		t.Errorf("first invoice status: want 'Paga', got %q", resp.Response[0].Status)
	}
	if resp.Response[1].Status != "Em aberto" {
		t.Errorf("second invoice status: want 'Em aberto', got %q", resp.Response[1].Status)
	}
}

func TestGetOpenInvoicesByContract_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "getcontractbillets/12345") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": []map[string]any{
				{
					"id":             123123,
					"title":          "FAT123456789",
					"expirationDate": "2024-05-30T00:00:00",
					"parcel":         1,
					"typefulLine":    "00000.00000",
					"link":           "https://example.com/boleto/abc123",
					"pixQRCode":      "pix-data",
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.GetOpenInvoicesByContract(context.Background(), 12345)
	if err != nil {
		t.Fatalf("GetOpenInvoicesByContract: %v", err)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("expected 1 billet, got %d", len(resp.Response))
	}
	billet := resp.Response[0]
	if billet.ID != 123123 {
		t.Errorf("id: want 123123, got %d", billet.ID)
	}
	if billet.Title != "FAT123456789" {
		t.Errorf("title: want 'FAT123456789', got %q", billet.Title)
	}
}

func TestRegisterPayment_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"synGwTransactionId": "a1387044-7154-41df-9e72-ec75c842126b",
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.RegisterPayment(context.Background(), RegisterPaymentRequest{
		TransactionID:              "550e8400-e29b-41d4-a716-446655440000",
		FinancialReceivableTitleID: 123456,
		PaidAmount:                 99.90,
		Message:                    "Pagamento via PIX",
		BankAccountCode:            "001",
		PaymentFormCode:            "PIX",
		ReceiptDate:                "2024-05-01",
		ClientPaidDate:             "2024-05-01",
	})
	if err != nil {
		t.Fatalf("RegisterPayment: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Response.SynGwTransactionID != "a1387044-7154-41df-9e72-ec75c842126b" {
		t.Errorf("unexpected synGwTransactionId: %q", resp.Response.SynGwTransactionID)
	}
}

func TestGeneratePix_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("receivableId") != "123456" {
			t.Errorf("receivableId: want '123456', got %q", r.URL.Query().Get("receivableId"))
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"registerId":     11,
				"qrCode":         "00020101021226870014br.gov.bcb.pix...",
				"totalAmount":    101.03,
				"interestAmount": 1.17,
				"fineAmount":     1.96,
				"transactionId":  "txn-uuid-here",
				"hybridBillet":   false,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.GeneratePix(context.Background(), 123456)
	if err != nil {
		t.Fatalf("GeneratePix: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Response.RegisterID != 11 {
		t.Errorf("registerId: want 11, got %d", resp.Response.RegisterID)
	}
	if resp.Response.TotalAmount != 101.03 {
		t.Errorf("totalAmount: want 101.03, got %f", resp.Response.TotalAmount)
	}
	if resp.Response.HybridBillet {
		t.Error("expected hybridBillet=false")
	}
}

func TestConsultPaymentStatus_Success(t *testing.T) {
	const txnID = "1eb56ef0-dfa1-11eb-8d19-0242ac130003"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("syngwtransactionid") != txnID {
			t.Errorf("syngwtransactionid: want %q, got %q", txnID, r.URL.Query().Get("syngwtransactionid"))
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"response": map[string]any{
				"integratorTransactionId": txnID,
				"status":                  map[string]any{"value": 1, "label": "Processado com sucesso"},
				"processed":               "2024-01-01T16:11:52",
				"message":                 "OK",
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	resp, err := client.ConsultPaymentStatus(context.Background(), txnID)
	if err != nil {
		t.Fatalf("ConsultPaymentStatus: %v", err)
	}
	if resp.Response.IntegratorTransactionID != txnID {
		t.Errorf("integratorTransactionId: want %q, got %q", txnID, resp.Response.IntegratorTransactionID)
	}
	if resp.Response.Status.Label != "Processado com sucesso" {
		t.Errorf("status.label: want 'Processado com sucesso', got %q", resp.Response.Status.Label)
	}
}

func TestGetInvoicePDF_Download(t *testing.T) {
	pdfContent := []byte("%PDF-1.4 test content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "GetBillet/42") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write(pdfContent)
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	data, ct, err := client.GetInvoicePDF(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetInvoicePDF: %v", err)
	}
	if ct != "application/pdf" {
		t.Errorf("content-type: want 'application/pdf', got %q", ct)
	}
	if string(data) != string(pdfContent) {
		t.Errorf("data mismatch")
	}
}

func TestRegisterRenegotiation_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"parcels": []map[string]any{
				{"id": 1, "title": "FAT001", "expirationDate": "2024-12-01", "amount": 150.00},
				{"id": 2, "title": "FAT002", "expirationDate": "2025-01-01", "amount": 150.00},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	result, err := client.RegisterRenegotiation(context.Background(), RegisterRenegotiationRequest{
		ReceivableIDs: []int{10, 11, 12},
		Discount:      0,
		Fine:          5.0,
		Interest:      2.5,
		Observation:   "Renegociação aprovada",
		Parcels: []RenegotiationParcel{
			{Number: 1, ExpirationDate: "2024-12-01"},
			{Number: 2, ExpirationDate: "2025-01-01"},
		},
	})
	if err != nil {
		t.Fatalf("RegisterRenegotiation: %v", err)
	}
	if len(result.Parcels) != 2 {
		t.Fatalf("expected 2 parcels, got %d", len(result.Parcels))
	}
	if result.Parcels[0].Title != "FAT001" {
		t.Errorf("first parcel title: want 'FAT001', got %q", result.Parcels[0].Title)
	}
}
