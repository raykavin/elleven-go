package elleven

import (
	"context"
	"fmt"
	"net/http"
)

// Billing (Faturamento)

// ContractClientInfo holds the client identification fields for contract creation.
type ContractClientInfo struct {
	TxID string `json:"txId"`
}

// ContractInfo holds the core contract configuration fields.
type ContractInfo struct {
	ContractType                string `json:"contractType"`
	CompanyPlaceTxID            string `json:"companyPlaceTxId"`
	CRMCampaignCode             string `json:"crmCampaignCode"`
	CRMPriceListCode            string `json:"crmPriceListCode"`
	PaymentFormCode             string `json:"paymentFormCode"`
	FinancialCollectionTypeCode string `json:"financialCollectionTypeCode"`
	DateManagementCode          string `json:"dateManagementCode"`
	CollectionDay               int    `json:"collectionDay"`
}

// ContractServiceInfo holds a service product for contract creation.
type ContractServiceInfo struct {
	Code            string  `json:"code"`
	Quantity        int     `json:"quantity"`
	Price           float64 `json:"price"`
	PaymentFormCode string  `json:"paymentFormCode"`
}

// ContractActivationInfo holds activation configuration for a new contract.
type ContractActivationInfo struct {
	IncidentTypeCode            string `json:"incidentTypeCode"`
	GroupInMonthlyFee           bool   `json:"groupInMonthlyFee"`
	FirstParcelDate             string `json:"firstParcelDate"` // YYYY-MM-DD
	PaymentFormCode             string `json:"paymentFormCode"`
	FinancialCollectionTypeCode string `json:"financialCollectionTypeCode"`
	PaymentConditionCode        string `json:"paymentConditionCode"`
}

// ContractIntegratorInfo holds integrator metadata for contract creation.
type ContractIntegratorInfo struct {
	IntegratorAlias     string `json:"integratorAlias"`
	IntegrationCode     string `json:"integrationCode"`
	IntegrationTeamCode string `json:"integrationTeamCode"`
}

// CreateContractRequest is the request body for creating a contract directly via billing.
type CreateContractRequest struct {
	ClientInformation     ContractClientInfo     `json:"clientInformation"`
	ContractInformation   ContractInfo           `json:"contractInformation"`
	ServicesInformation   []ContractServiceInfo  `json:"servicesInformation"`
	ActivationInformation ContractActivationInfo `json:"activationInformation"`
	IntegratorInformation ContractIntegratorInfo `json:"integratorInformation"`
}

// CreateContractResult is the response payload when a contract is created.
type CreateContractResult struct {
	ContractNumber       string `json:"contractNumber"`
	ActivationAssignment int    `json:"activationAssignment"`
}

// CreateContractResponse wraps the contract creation response.
type CreateContractResponse struct {
	Success          bool                 `json:"Success"`
	Message          *string              `json:"message"`
	Response         CreateContractResult `json:"response"`
	DataResponseType string               `json:"dataResponseType"`
	ElapsedTime      *float64             `json:"elapsedTime"`
}

// SaleOrderItem represents a single line item in a sale order.
type SaleOrderItem struct {
	ID        string  `json:"id"`
	Quantity  int     `json:"quantity"`
	UnitValue float64 `json:"unitValue"`
}

// GenerateSaleOrderRequest is the request body for generating a sale order.
type GenerateSaleOrderRequest struct {
	CompanyPlaceTxID          string          `json:"companyPlaceTxId"`
	ClientTxID                string          `json:"clientTxId"`
	ContractNumber            string          `json:"contractNumber"`
	PriceListID               string          `json:"priceListId"`
	FinancialCollectionTypeID string          `json:"financialCollectionTypeId"`
	PaymentConditionID        string          `json:"paymentConditionId"`
	Items                     []SaleOrderItem `json:"items"`
}

// EventualValueType defines the type of an eventual (one-time) billing value.
type EventualValueType int

const (
	EventualValueTypeCredit EventualValueType = 0
	EventualValueTypeDebit  EventualValueType = 1
)

// CreateEventualValueRequest is the request body for creating an eventual billing value.
type CreateEventualValueRequest struct {
	Type                           EventualValueType `json:"type"`
	ContractID                     int               `json:"contractId"`
	ContractItemID                 int               `json:"contracItemtId"` 
	ContractConfigurationBillingID int               `json:"contractConfigurationBillingId"`
	ServiceProductCode             string            `json:"serviceProductCode"`
	MonthYearType                  int               `json:"monthYearType"`
	MonthYear                      string            `json:"monthYear"`        // YYYY-MM-DD
	InitialMonthYear               string            `json:"initialMonthYear"` // YYYY-MM-DD
	FinalMonthYear                 string            `json:"finalMonthYear"`   // YYYY-MM-DD
	Description                    string            `json:"description"`
	PresentInvoiceNote             bool              `json:"presentInvoiceNote"`
	UnitAmount                     float64           `json:"unitAmount"`
	Units                          int               `json:"units"`
}

// ContractDocument represents a document file associated with a contract.
type ContractDocument struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	DocumentationType string `json:"documentationType"`
	Validity          struct {
		Begin *string `json:"begin"`
		Final *string `json:"final"`
	} `json:"validity"`
	File string `json:"file"` // URL to the file
}

// ApproveContractRequest is the request body for approving a contract.
type ApproveContractRequest struct {
	ProportionalityType           int    `json:"proportionalityType"`
	ChangeBeginningDateOnApproval bool   `json:"changeBeginningDateOnApproval"`
	ChangeLoyaltyDateOnApproval   bool   `json:"changeLoyaltyDateOnApproval"`
	ApprovalDate                  string `json:"approvalDate"` // RFC3339
}

// ProportionalityItem represents a proportionality entry from a contract approval.
type ProportionalityItem struct {
	Description        string  `json:"description"`
	Amount             float64 `json:"amount"`
	Competence         string  `json:"competence"`
	ServiceDescription string  `json:"serviceDescription"`
	Type               int     `json:"type"`
	CompetenceType     int     `json:"competenceType"`
}

// ApproveContractData is the response data for contract approval.
type ApproveContractData struct {
	PersonUsers       []any                 `json:"personUsers"`
	Proportionalities []ProportionalityItem `json:"proportionalities"`
}

// ApproveContractResponse wraps the contract approval response.
type ApproveContractResponse struct {
	Success     bool                `json:"success"`
	Data        ApproveContractData `json:"data"`
	ElapsedTime float64             `json:"elapsedTime"`
}

// CreateContract creates a new contract directly via the billing module.
//
// POST /external/billing/contracts/create
func (c *Client) CreateContract(ctx context.Context, req CreateContractRequest) (*CreateContractResponse, error) {
	var result CreateContractResponse
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/billing/contracts/create"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GenerateSaleOrder generates a sale order for an existing contract.
//
// POST /external/integrations/thirdparty/salerequest
func (c *Client) GenerateSaleOrder(ctx context.Context, req GenerateSaleOrderRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/salerequest"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateEventualValue adds a one-time credit or debit value to a contract billing period.
//
// POST /external/integrations/thirdparty/contract/contracteventualvalues
func (c *Client) CreateEventualValue(ctx context.Context, req CreateEventualValueRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/contract/contracteventualvalues"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDocumentsByContract lists documents attached to a contract filtered by documentation type.
//
// GET /external/integrations/thirdparty/getfilesbycontract/{contractNumber}/documentationtype/{documentationTypeID}
func (c *Client) ListDocumentsByContract(ctx context.Context, contractNumber string, documentationTypeID int) (*APIResponse[[]ContractDocument], error) {
	u := c.apiURL(fmt.Sprintf(
		"/external/integrations/thirdparty/getfilesbycontract/%s/documentationtype/%d",
		contractNumber, documentationTypeID,
	))

	var result APIResponse[[]ContractDocument]
	if err := c.doJSON(ctx, http.MethodGet, u, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadAttachmentToContract uploads a file attachment to a contract.
//
// POST /external/integrations/thirdparty/contract/contractuploads
// Query params: contractNumber, documentationTypeCode
func (c *Client) UploadAttachmentToContract(
	ctx context.Context,
	contractNumber string,
	documentationTypeCode string,
	fileName string,
	fileData []byte,
) (*APIResponse[any], error) {
	u := buildURL(
		c.apiURL("/external/integrations/thirdparty/contract/contractuploads"),
		map[string]string{
			"contractNumber":        contractNumber,
			"documentationTypeCode": documentationTypeCode,
		},
	)

	var result APIResponse[any]
	if err := c.doMultipart(ctx, u, "File", fileName, fileData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UnlockContract unlocks a suspended/blocked contract.
//
// POST /external/integrations/thirdparty/contracts/unlock/{contractNumber}
func (c *Client) UnlockContract(ctx context.Context, contractNumber string) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/contracts/unlock/%s", contractNumber)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadContractAttachment downloads a file attached to a contract.
//
// Returns raw file bytes and Content-Type header value.
//
// GET /external/integrations/thirdparty/contract/contractuploads/getfile
func (c *Client) DownloadContractAttachment(ctx context.Context, attachmentID, contractID int) ([]byte, string, error) {
	u := buildURL(
		c.apiURL("/external/integrations/thirdparty/contract/contractuploads/getfile"),
		map[string]string{
			"id":         fmt.Sprintf("%d", attachmentID),
			"contractId": fmt.Sprintf("%d", contractID),
		},
	)
	return c.doDownload(ctx, u)
}

// ApproveContract approves a pending contract, optionally generating proportionalities.
//
// PUT /external/integrations/thirdparty/contracts/approve/{contractID}
func (c *Client) ApproveContract(ctx context.Context, contractID int, req ApproveContractRequest) (*ApproveContractResponse, error) {
	var result ApproveContractResponse
	if err := c.doJSON(ctx, http.MethodPut,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/contracts/approve/%d", contractID)),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadDANFE downloads a DANFE (Nota Fiscal Eletrônica) PDF for a billing document.
//
// Returns raw PDF bytes and Content-Type header value.
//
// GET /external/integrations/thirdparty/billings/invoices/download/{documentID}
func (c *Client) DownloadDANFE(ctx context.Context, documentID int) ([]byte, string, error) {
	u := c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/billings/invoices/download/%d", documentID))
	return c.doDownload(ctx, u)
}
