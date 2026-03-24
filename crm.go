package elleven

import (
	"context"
	"net/http"
)

// CRM

// LeadPersonalData holds personal identification fields for a lead.
type LeadPersonalData struct {
	Name        string `json:"name"`
	Name2       string `json:"name2"`
	TypeTxID    int    `json:"typeTxId"`
	TxID        string `json:"txId"`
	BirthDate   string `json:"birthDate,omitempty"` // RFC3339
	Identity    string `json:"identity"`
	CivilStatus *int   `json:"civilStatus"`
	Gender      *int   `json:"gender"`
	ParentsName string `json:"parentsName"`
}

// LeadContact holds contact information for a lead.
type LeadContact struct {
	Phone                string `json:"phone"`
	CellPhone            string `json:"cellPhone"`
	CellPhoneHasWhatsApp bool   `json:"cellPhoneHasWhatsApp"`
	CellPhoneHasTelegram bool   `json:"cellPhoneHasTelegram"`
	Email                string `json:"email"`
	Website              string `json:"website"`
}

// LeadAddress holds address fields for a lead.
type LeadAddress struct {
	PostalCode        string `json:"postalCode"`
	Street            string `json:"street"`
	StreetType        string `json:"streetType"`
	Number            string `json:"number"`
	Neighborhood      string `json:"neighborhood"`
	AddressComplement string `json:"addressComplement"`
	AddressReference  string `json:"addressReference"`
	CodeCityID        int    `json:"codeCityId"`
	City              string `json:"city"`
	State             string `json:"state"`
	Country           string `json:"country"`
}

// LeadIntegratorData holds integration metadata for a lead creation request.
type LeadIntegratorData struct {
	CRMContactOriginCode  string `json:"crmContactOriginCode"`
	CRMFormCode           string `json:"crmFormCode"`
	IntegratorAlias       string `json:"integratorAlias"`
	IntegrationCode       string `json:"integrationCode"`
	IntegrationTeamCode   string `json:"integrationTeamCode"`
	IntegrationSellerTxID string `json:"integrationSellerTxId"`
}

// CreateLeadRequest is the request body for creating a CRM lead.
type CreateLeadRequest struct {
	PersonalData   LeadPersonalData   `json:"personalData"`
	Contact        LeadContact        `json:"contact"`
	Address        LeadAddress        `json:"address"`
	IntegratorData LeadIntegratorData `json:"integratorData"`
}

// SaleServicePhoneInfo contains phone number details for a service product.
type SaleServicePhoneInfo struct {
	PhoneNumber     string `json:"PhoneNumber"`
	PhoneNumberType string `json:"PhoneNumberType"`
	IsPortability   bool   `json:"IsPortability"`
}

// SaleServiceProduct describes a service product to be included in a sale.
type SaleServiceProduct struct {
	Code                    string                 `json:"code"`
	Quantity                int                    `json:"quantity"`
	Amount                  float64                `json:"amount"`
	PaymentFormCode         string                 `json:"paymentFormCode"`
	PhoneGroupID            string                 `json:"phoneGroupId"`
	PhoneNumberInformations []SaleServicePhoneInfo `json:"PhoneNumberInformations,omitempty"`
}

// StartSaleRequest is the request body for opening a new negotiation (sale).
//
// This is one of the most complex endpoints in the API. Fields not shown in the
// documentation as structured objects are modelled as map[string]any to
// allow flexibility until the API provides a full schema.
type StartSaleRequest struct {
	TxID                                  string               `json:"txId"`
	CompanyPlaceTxID                      string               `json:"companyPlaceTxId"`
	ThirdSellerTxID                       string               `json:"thirdSellerTxId"`
	ResponsibleSellerTxID                 string               `json:"responsibleSellerTxId"`
	TeamCode                              string               `json:"TeamCode"`
	ContractType                          string               `json:"contractType"`
	CRMCampaignCode                       string               `json:"crmCampaignCode"`
	CRMPriceListCode                      string               `json:"crmPriceListCode"`
	ContractPaymentFormCode               string               `json:"contractPaymentFormCode"`
	ContractFinancialCollectionTypeCode   string               `json:"contractFinancialCollectionTypeCode"`
	ContractDateManagementCode            string               `json:"contractDateManagementCode"`
	ContractCollectionDay                 string               `json:"contractCollectionDay"`
	ActivationIncidentTypeCode            string               `json:"activationIncidentTypeCode"`
	ActivationGroupInMonthlyFee           bool                 `json:"activationGroupInMonthlyFee"`
	ActivationFirstParcelDate             string               `json:"activationFirstParcelDate"` // YYYY-MM-DD
	ActivationPaymentFormCode             string               `json:"activationPaymentFormCode"`
	ActivationFinancialCollectionTypeCode string               `json:"activationFinancialCollectionTypeCode"`
	ActivationPaymentConditionCode        string               `json:"activationPaymentConditionCode"`
	ServiceProducts                       []SaleServiceProduct `json:"serviceProducts"`

	// ClientSaleInfo, ActivationInformation, IntegratorInformation, GainInformation
	// are flexible objects. Pass nil or a custom struct as needed.
	ClientSaleInfo        any    `json:"clientSaleInfo,omitempty"`
	ActivationInformation any    `json:"activationInformation,omitempty"`
	IntegratorInformation any    `json:"integratorInformation,omitempty"`
	GainInformation       any    `json:"gainInformation,omitempty"`
	Contacts              []any  `json:"contacts,omitempty"`
	Observations          string `json:"observations"`
}

// StartSaleResult is the response payload for a new sale/negotiation.
type StartSaleResult struct {
	ProtocolID int `json:"protocolId"`
}

// ContractServicePhoneInfo contains phone info for a contract service addition.
type ContractServicePhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PhoneNumberType int    `json:"phoneNumberType"`
	IsPortability   bool   `json:"isPortability"`
}

// NewContractService describes a service to be added to a contract.
type NewContractService struct {
	Code                          string                     `json:"code"`
	Quantity                      int                        `json:"quantity"`
	Price                         float64                    `json:"price"`
	PaymentFormCode               string                     `json:"paymentFormCode"`
	ServiceTagOption              int                        `json:"serviceTagOption"`
	ContractServiceTag            string                     `json:"contractServiceTag"`
	ContractServiceTagDescription string                     `json:"contractServiceTagDescription"`
	PhoneGroupID                  string                     `json:"phoneGroupId"`
	PhoneNumberInformations       []ContractServicePhoneInfo `json:"phoneNumberInformations,omitempty"`
}

// AddContractServicesRequest is the request body for adding services to a contract.
type AddContractServicesRequest struct {
	ContractNumber          string               `json:"contractNumber"`
	CRMCampaignCode         string               `json:"crmCampaignCode"`
	CRMPriceListCode        string               `json:"crmPriceListCode"`
	GenerateProportionality bool                 `json:"generateProportionality"`
	NewServices             []NewContractService `json:"newServices"`
}

// ContractServiceModification extends NewContractService with the existing item ID.
type ContractServiceModification struct {
	OldContractItemID int `json:"oldContractItemId"`
	NewContractService
}

// ChangeContractServicesRequest is the request body for modifying contract services.
type ChangeContractServicesRequest struct {
	ContractNumber          string                        `json:"contractNumber"`
	CRMCampaignCode         string                        `json:"crmCampaignCode"`
	CRMPriceListCode        string                        `json:"crmPriceListCode"`
	GenerateProportionality bool                          `json:"generateProportionality"`
	ChangedServices         []ContractServiceModification `json:"changedServices"`
}

// RemoveContractItem identifies a contract item by its ID for removal.
type RemoveContractItem struct {
	ContractItem int `json:"contractItem"`
}

// RemoveContractServicesRequest is the request body for removing services from a contract.
type RemoveContractServicesRequest struct {
	ContractNumber string               `json:"contractNumber"`
	RemoveService  []RemoveContractItem `json:"removeService"`
}

// ContractServiceOperationResult is the response for contract service operations.
// Note: the API returns inconsistent casing (Sucess/True) mapped faithfully.
type ContractServiceOperationResult struct {
	Success bool   `json:"Sucess"` // Note: API uses "Sucess" (typo in API)
	Message string `json:"Message"`
}

// ViabilityAddress holds address fields for a viability check.
type ViabilityAddress struct {
	Address      string `json:"address"`
	Number       string `json:"number"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
}

// ViabilityPerson holds person reference data for a viability check.
type ViabilityPerson struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TxIDType string `json:"txIdType"`
	TxID     string `json:"txId"`
	Phone    string `json:"phone"`
}

// ViabilityCampaign holds campaign reference data for a viability check.
type ViabilityCampaign struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// VerifyViabilityRequest is the request body for checking address coverage viability.
type VerifyViabilityRequest struct {
	FullAddress ViabilityAddress  `json:"fullAddress"`
	Distance    string            `json:"distance"`
	Lead        ViabilityPerson   `json:"lead"`
	Seller      ViabilityPerson   `json:"seller"`
	Campaign    ViabilityCampaign `json:"campaign"`
}

// ViabilityResult is the response payload for a viability check.
type ViabilityResult struct {
	Viability bool `json:"viability"`
	CTOs      int  `json:"ctos"`
	Ports     int  `json:"ports"`
}

// CancelSaleRequest is the request body for cancelling a negotiation.
type CancelSaleRequest struct {
	ProtocolID  int    `json:"protocolId"`
	Description string `json:"description"`
}

// ContractTypeService represents a service product inside a contract type.
type ContractTypeService struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

// ContractType represents a contract type with its associated service products.
type ContractType struct {
	Code                        string                `json:"code"`
	Title                       string                `json:"title"`
	Description                 string                `json:"description"`
	CollectionDays              []int                 `json:"collectionDays"`
	ContractTypesServiceProduct []ContractTypeService `json:"contractTypesServiceProduct"`
}

// CampaignPriceListProduct represents a product inside a campaign price list.
type CampaignPriceListProduct struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
}

// CampaignPriceList represents a price list within a campaign.
type CampaignPriceList struct {
	Code        string                     `json:"code"`
	Title       string                     `json:"title"`
	Description *string                    `json:"description"`
	Products    []CampaignPriceListProduct `json:"products"`
}

// CampaignCity represents a city where a campaign is available.
type CampaignCity struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CampaignRegion represents a region covered by a campaign.
type CampaignRegion struct {
	Code   string         `json:"code"`
	Name   string         `json:"name"`
	Cities []CampaignCity `json:"cities"`
}

// Campaign represents a CRM campaign with price lists and regions.
type Campaign struct {
	Code        string              `json:"code"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Regions     []CampaignRegion    `json:"regions"`
	PriceLists  []CampaignPriceList `json:"priceLists"`
}

// CreateLead creates a new CRM lead.
//
// POST /external/crm/leads/create
func (c *Client) CreateLead(ctx context.Context, req CreateLeadRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/crm/leads/create"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// StartSale opens a new CRM negotiation (sale) and returns a protocol ID.
//
// POST /external/integrations/thirdparty/crm/startsale
func (c *Client) StartSale(ctx context.Context, req StartSaleRequest) (*APIResponse[StartSaleResult], error) {
	var result APIResponse[StartSaleResult]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/crm/startsale"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddContractServices adds one or more service products to an existing contract.
//
// POST /external/integrations/thirdparty/contract/addcontractservices
func (c *Client) AddContractServices(ctx context.Context, req AddContractServicesRequest) (*ContractServiceOperationResult, error) {
	var result ContractServiceOperationResult
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/contract/addcontractservices"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ChangeContractServices modifies existing service products on a contract.
//
// POST /external/integrations/thirdparty/contract/changecontractservices
func (c *Client) ChangeContractServices(ctx context.Context, req ChangeContractServicesRequest) (*ContractServiceOperationResult, error) {
	var result ContractServiceOperationResult
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/contract/changecontractservices"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveContractServices removes service items from an existing contract.
//
// POST /external/integrations/thirdparty/contract/removecontractitems
func (c *Client) RemoveContractServices(ctx context.Context, req RemoveContractServicesRequest) (*ContractServiceOperationResult, error) {
	var result ContractServiceOperationResult
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/contract/removecontractitems"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyViability checks whether a given address has service coverage.
//
// POST /external/integrations/thirdparty/verifyviability
func (c *Client) VerifyViability(ctx context.Context, req VerifyViabilityRequest) (*APIResponse[ViabilityResult], error) {
	var result APIResponse[ViabilityResult]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/verifyviability"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelSale cancels an existing CRM negotiation by protocol ID.
//
// POST /external/integrations/thirdparty/crm/cancelsale
func (c *Client) CancelSale(ctx context.Context, req CancelSaleRequest) (*APIResponse[string], error) {
	var result APIResponse[string]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/crm/cancelsale"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetContractTypesAndServices returns available contract types with their service products.
//
// GET /external/integrations/thirdparty/crm/contracttypesandservices
func (c *Client) GetContractTypesAndServices(ctx context.Context) (*APIResponse[[]ContractType], error) {
	var result APIResponse[[]ContractType]
	if err := c.doJSON(ctx, http.MethodGet,
		c.apiURL("/external/integrations/thirdparty/crm/contracttypesandservices"),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCampaignsAndPriceListServices returns CRM campaigns with their price lists and products.
//
// GET /external/integrations/thirdparty/crm/campaignsandpricelistservices
func (c *Client) GetCampaignsAndPriceListServices(ctx context.Context) (*APIResponse[[]Campaign], error) {
	var result APIResponse[[]Campaign]
	if err := c.doJSON(ctx, http.MethodGet,
		c.apiURL("/external/integrations/thirdparty/crm/campaignsandpricelistservices"),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
