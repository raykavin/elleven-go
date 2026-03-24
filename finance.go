package elleven

import (
	"context"
	"fmt"
	"net/http"
)

// Finance (Financeiro)

// InvoiceAmount holds the detailed amounts for an invoice.
type InvoiceAmount struct {
	Value      float64 `json:"value"`
	FinalValue float64 `json:"finalValue"`
	Discount   float64 `json:"discount"`
	Fine       float64 `json:"fine"`
	Interest   float64 `json:"interest"`
}

// InvoiceBillet holds billet (boleto) and PIX data for an invoice.
type InvoiceBillet struct {
	BankTitleNumber *string       `json:"bankTitleNumber"`
	Barcode         *string       `json:"barcode"`
	TypefulLine     string        `json:"typefulLine"`
	PixQRCode       string        `json:"pixQRCode"`
	Title           string        `json:"title"`
	IssueDate       string        `json:"issueDate"`
	ExpirationDate  string        `json:"expirationDate"`
	ProcessingDate  string        `json:"processingDate"`
	Amount          InvoiceAmount `json:"amount"`
}

// Invoice represents a financial receivable title (invoice/fatura).
type Invoice struct {
	ID             int           `json:"id"`
	Billet         InvoiceBillet `json:"billet"`
	Client         any           `json:"client"`
	CompanyPlace   any           `json:"companyPlace"`
	Bank           any           `json:"bank"`
	CollectionType any           `json:"collectionType"`

	// Status is only present on GetAllInvoicesByTxID (e.g. "Em aberto", "Paga", "Cancelada", "Vencida")
	Status string `json:"status,omitempty"`
}

// ContractBillet represents a simplified billet as returned by GetOpenInvoicesByContract.
type ContractBillet struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	ExpirationDate string `json:"expirationDate"`
	Parcel         int    `json:"parcel"`
	TypefulLine    string `json:"typefulLine"`
	Link           string `json:"link"`
	PixQRCode      string `json:"pixQRCode"`
}

// PaymentStatus holds the status label and value for a payment.
type PaymentStatus struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// PaymentStatusResult is the response payload for ConsultPaymentStatus.
type PaymentStatusResult struct {
	IntegratorTransactionID string        `json:"integratorTransactionId"`
	Status                  PaymentStatus `json:"status"`
	Processed               string        `json:"processed"`
	Message                 string        `json:"message"`
}

// RegisterPaymentRequest is the request body for settling an invoice payment.
type RegisterPaymentRequest struct {
	// TransactionID is a unique GUID for this payment transaction (idempotency key).
	TransactionID              string  `json:"transactionId"`
	FinancialReceivableTitleID int     `json:"financialReceivableTitleId"`
	PaidAmount                 float64 `json:"paidAmount"`
	Message                    string  `json:"message"`
	BankAccountCode            string  `json:"BankAccountCode"`
	PaymentFormCode            string  `json:"PaymentFormCode"`
	ReceiptDate                string  `json:"ReceiptDate"`    // YYYY-MM-DD
	ClientPaidDate             string  `json:"ClientPaidDate"` // YYYY-MM-DD
}

// RegisterPaymentResult is the response payload for a successful payment registration.
type RegisterPaymentResult struct {
	SynGwTransactionID string `json:"synGwTransactionId"`
}

// GeneratePixRequest identifies the invoice for PIX generation.
// Passed as a query parameter, not a body.
type GeneratePixResult struct {
	RegisterID     int     `json:"registerId"`
	QRCode         string  `json:"qrCode"`
	TotalAmount    float64 `json:"totalAmount"`
	InterestAmount float64 `json:"interestAmount"`
	FineAmount     float64 `json:"fineAmount"`
	TransactionID  string  `json:"transactionId"`
	HybridBillet   bool    `json:"hybridBillet"`
}

// RenegotiationInfoRequest is the request body for getting renegotiation information.
type RenegotiationInfoRequest struct {
	ReceivableIDs []int  `json:"receivableIds"`
	Date          string `json:"date"` // YYYY-MM-DD
}

// RenegotiationReceivable represents a single receivable in a renegotiation.
type RenegotiationReceivable struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Fine     float64 `json:"fine"`
	Interest float64 `json:"interest"`
}

// RenegotiationInfoResult is the response payload for renegotiation information.
type RenegotiationInfoResult struct {
	TotalFine     float64                   `json:"totalFine"`
	TotalInterest float64                   `json:"totalInterest"`
	Receivables   []RenegotiationReceivable `json:"receivables"`
	Date          string                    `json:"date"`
}

// RenegotiationParcel describes a single installment in a renegotiation.
type RenegotiationParcel struct {
	Number         int      `json:"number"`
	Amount         *float64 `json:"amount"`         // Nullable: API allows null for calculated amounts
	ExpirationDate string   `json:"expirationDate"` // YYYY-MM-DD
}

// RegisterRenegotiationRequest is the request body for registering a debt renegotiation.
type RegisterRenegotiationRequest struct {
	ReceivableIDs               []int                 `json:"receivableIds"`
	Discount                    float64               `json:"discount"`
	Fine                        float64               `json:"fine"`
	Interest                    float64               `json:"interest"`
	FinancialCollectionTypeCode any                   `json:"financialCollectionTypeCode"` // int or string per API
	Observation                 string                `json:"observation"`
	BillingInvoiceAntecipated   bool                  `json:"billingInvoiceAntecipad"`
	Parcels                     []RenegotiationParcel `json:"parcels"`
}

// RenegotiationResultParcel is a single parcel in the renegotiation result.
type RenegotiationResultParcel struct {
	ID             int     `json:"id"`
	Title          string  `json:"title"`
	ExpirationDate string  `json:"expirationDate"`
	Amount         float64 `json:"amount"`
}

// RegisterRenegotiationResult is the response payload after a renegotiation is registered.
type RegisterRenegotiationResult struct {
	Parcels []RenegotiationResultParcel `json:"parcels"`
}

// GetOpenInvoicesByTxID returns open (unpaid) invoices for a CPF/CNPJ.
//
// GET /external/integrations/thirdparty/getopentitlesbytxid/{txId}
func (c *Client) GetOpenInvoicesByTxID(ctx context.Context, txID string) (*APIResponse[[]Invoice], error) {
	var result APIResponse[[]Invoice]
	if err := c.doJSON(ctx, http.MethodGet,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/getopentitlesbytxid/%s", txID)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllInvoicesByTxID returns all invoices (any status) for a CPF/CNPJ.
//
// GET /external/integrations/thirdparty/gettitlesbytxid/{txId}
func (c *Client) GetAllInvoicesByTxID(ctx context.Context, txID string) (*APIResponse[[]Invoice], error) {
	var result APIResponse[[]Invoice]
	if err := c.doJSON(ctx, http.MethodGet,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/gettitlesbytxid/%s", txID)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOpenInvoicesByContract returns open invoices for a specific contract.
//
// GET /external/integrations/thirdparty/getcontractbillets/{contractID}
func (c *Client) GetOpenInvoicesByContract(ctx context.Context, contractID int) (*APIResponse[[]ContractBillet], error) {
	var result APIResponse[[]ContractBillet]
	if err := c.doJSON(ctx, http.MethodGet,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/getcontractbillets/%d", contractID)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInvoicePDF downloads the printed billet (boleto) PDF for an invoice.
//
// Returns raw PDF bytes and Content-Type header value.
//
// GET /external/integrations/thirdparty/GetBillet/{invoiceID}
func (c *Client) GetInvoicePDF(ctx context.Context, invoiceID int) ([]byte, string, error) {
	u := c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/GetBillet/%d", invoiceID))
	return c.doDownload(ctx, u)
}

// ConsultPaymentStatus checks the processing status of a payment by its SynGW transaction ID.
//
// GET /external/integrations/thirdparty/consultpayment?syngwtransactionid={uuid}
func (c *Client) ConsultPaymentStatus(ctx context.Context, synGWTransactionID string) (*APIResponse[PaymentStatusResult], error) {
	u := buildURL(
		c.apiURL("/external/integrations/thirdparty/consultpayment"),
		map[string]string{"syngwtransactionid": synGWTransactionID},
	)

	var result APIResponse[PaymentStatusResult]
	if err := c.doJSON(ctx, http.MethodGet, u, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RegisterPayment settles an invoice by registering a payment (receipt).
//
// POST /external/integrations/thirdparty/receivepayment
func (c *Client) RegisterPayment(ctx context.Context, req RegisterPaymentRequest) (*APIResponse[RegisterPaymentResult], error) {
	var result APIResponse[RegisterPaymentResult]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/receivepayment"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GeneratePix generates a PIX QR code for an existing invoice (receivable).
//
// POST /external/integrations/thirdparty/billings/registerpix?receivableId={id}
func (c *Client) GeneratePix(ctx context.Context, receivableID int) (*APIResponse[GeneratePixResult], error) {
	u := buildURL(
		c.apiURL("/external/integrations/thirdparty/billings/registerpix"),
		map[string]string{"receivableId": fmt.Sprintf("%d", receivableID)},
	)

	var result APIResponse[GeneratePixResult]
	if err := c.doJSON(ctx, http.MethodPost, u, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRenegotiationInfo retrieves fine and interest information for a set of receivables
// to support building a renegotiation proposal.
//
// GET /external/integrations/thirdparty/financial/getrenegotiationsinformations
// (body sent as JSON despite being a GET matches API documentation)
func (c *Client) GetRenegotiationInfo(ctx context.Context, req RenegotiationInfoRequest) (*APIResponse[RenegotiationInfoResult], error) {
	var result APIResponse[RenegotiationInfoResult]
	// NOTE: API documentation shows a GET with a JSON body. This is non-standard but
	// supported by many servers. We implement as-documented.
	if err := c.doJSON(ctx, http.MethodGet,
		c.apiURL("/external/integrations/thirdparty/financial/getrenegotiationsinformations"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RegisterRenegotiation creates a new debt renegotiation for a set of open receivables.
//
// POST /external/integrations/thirdparty/financial/registerrenegotiations
func (c *Client) RegisterRenegotiation(ctx context.Context, req RegisterRenegotiationRequest) (*RegisterRenegotiationResult, error) {
	var result RegisterRenegotiationResult
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/financial/registerrenegotiations"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
