package elleven

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// Suite (People / Customers)

// RegisterPersonRequest is the request body for registering a new person/customer.
type RegisterPersonRequest struct {
	TypeTxID          string `json:"typeTxId"`
	TxID              string `json:"txId"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Client            bool   `json:"client"`
	Situation         int    `json:"situation"`
	StreetType        string `json:"streetType"`
	PostalCode        string `json:"postalCode"`
	Street            string `json:"street"`
	Number            string `json:"number"`
	AddressComplement string `json:"addressComplement"`
	AddressReference  string `json:"addressReference"`
	Neighborhood      string `json:"neighborhood"`
	City              string `json:"city"`
	CodeCityID        int    `json:"codeCityId"`
	State             string `json:"state"`
	CodeCountry       string `json:"codeCountry"`
}

// RegisterPersonResult is the response payload when a person is registered.
type RegisterPersonResult struct{}

// UpdatePersonEmailRequest contains email fields for person update.
type UpdatePersonEmailRequest struct {
	Email    string `json:"email"`
	EmailNfe string `json:"emailNfe"`
}

// UpdatePersonPhoneRequest contains phone fields for person update.
type UpdatePersonPhoneRequest struct {
	Phone                string `json:"phone"`
	CellPhone            string `json:"cellPhone"`
	CellPhoneHasWhatsapp bool   `json:"cellPhoneHasWhatsapp"`
}

// UpdatePersonRequest is the request body for updating a person's contact info.
type UpdatePersonRequest struct {
	ID    int                      `json:"id"`
	Email UpdatePersonEmailRequest `json:"email,omitempty"`
	Phone UpdatePersonPhoneRequest `json:"phone,omitempty"`
}

// UpdatePersonResult is the response payload for person update.
type UpdatePersonResult struct {
	PeopleAddressID int    `json:"peopleAddressId"`
	Message         string `json:"message"`
}

// UpdateAddressRequest is the request body for updating a customer's address.
type UpdateAddressRequest struct {
	ID              int     `json:"id"`
	City            string  `json:"city"`
	Complement      string  `json:"complement"`
	IBGECode        string  `json:"ibgeCode"`
	Neighborhood    string  `json:"neighborhood"`
	Number          string  `json:"number"`
	PeopleAddressID int     `json:"peopleAddressId"`
	PublicPlace     string  `json:"publicPlace"`
	Reference       string  `json:"reference"`
	State           string  `json:"state"`
	ZipCode         string  `json:"zipCode"`
	PropertyType    *string `json:"propertyType"`
	Lat             string  `json:"lat"`
	Lng             string  `json:"lng"`
}

// UpdateAddressResult is the response payload for address update.
type UpdateAddressResult struct {
	PeopleAddressID int    `json:"peopleAddressId"`
	Message         string `json:"message"`
}

// Customer represents a full customer record returned by the list endpoint.
type Customer struct {
	ID                    any     `json:"id"` // API returns as string or number
	AddressComplement     string  `json:"addressComplement"`
	AddressReference      string  `json:"addressReference"`
	BirthDate             *string `json:"birthDate"`
	Candidate             bool    `json:"candidate"`
	Carrier               bool    `json:"carrier"`
	CellPhone1            string  `json:"cellPhone1"`
	City                  string  `json:"city"`
	CivilStatus           int     `json:"civilStatus"`
	Client                bool    `json:"client"`
	CodeCityID            int     `json:"codeCityId"`
	CodeCountry           string  `json:"codeCountry"`
	Collaborator          bool    `json:"collaborator"`
	Competitor            bool    `json:"competitor"`
	Country               string  `json:"country"`
	Email                 string  `json:"email"`
	EmailNfe              string  `json:"emailNfe"`
	Expert                bool    `json:"expert"`
	Gender                int     `json:"gender"`
	HouseLawyer           bool    `json:"houseLawyer"`
	Identity              string  `json:"identity"`
	MunicipalRegistration string  `json:"municipalRegistration"`
	Name                  string  `json:"name"`
	Name2                 string  `json:"name2"`
	Neighborhood          string  `json:"neighborhood"`
	Number                string  `json:"number"`
	ParentsName           string  `json:"parentsName"`
	PeopleAddressMainID   int     `json:"peopleAddressMainId"`
	Phone                 string  `json:"phone"`
	PostalCode            string  `json:"postalCode"`
	Situation             int     `json:"situation"`
	State                 string  `json:"state"`
	StateRegistration     string  `json:"stateRegistration"`
	StateRegistrationType int     `json:"stateRegistrationType"`
	Status                int     `json:"status"`
	Street                string  `json:"street"`
	StreetType            string  `json:"streetType"`
	TxID                  string  `json:"txId"`
	TxIDFormatted         string  `json:"txIdFormated"`
	TypeTxID              int     `json:"typeTxId"`
	UserID                int     `json:"userId"`
	Vendor                bool    `json:"vendor"`
}

// Location represents an ERP company location (place/branch).
type Location struct {
	ID            int    `json:"id"`
	Description   string `json:"description"`
	Code          string `json:"code"`
	TxID          string `json:"txId"`
	TxIDFormatted string `json:"txIdFormated"`
	City          string `json:"city"`
	Street        string `json:"street"`
}

// PersonAddress represents the main address of a person.
type PersonAddress struct {
	StreetType        string `json:"streetType"`
	Street            string `json:"street"`
	Number            string `json:"number"`
	AddressComplement string `json:"addressComplement"`
	Neighborhood      string `json:"neighborhood"`
	City              string `json:"city"`
	CodeCityID        int    `json:"codeCityId"`
	AddressReference  string `json:"addressReference"`
	State             string `json:"state"`
	PostalCode        string `json:"postalCode"`
	Longitude         string `json:"longitude"`
	Latitude          string `json:"latitude"`
}

// TitleAmount represents the financial amounts in a title/invoice.
type TitleAmount struct {
	Value      float64 `json:"value"`
	FinalValue float64 `json:"finalValue"`
	Discount   float64 `json:"discount"`
	Fine       float64 `json:"fine"`
	Interest   float64 `json:"interest"`
}

// TitleBillet represents a billing billet inside a title.
type TitleBillet struct {
	BankTitleNumber *string     `json:"bankTitleNumber"`
	Balance         float64     `json:"balance"`
	Title           string      `json:"title"`
	IssueDate       string      `json:"issueDate"`
	ExpirationDate  string      `json:"expirationDate"`
	ProcessingDate  string      `json:"processingDate"`
	Amount          TitleAmount `json:"amount"`
}

// PersonTitle represents a billing title associated with a person.
type PersonTitle struct {
	ID     int         `json:"id"`
	Billet TitleBillet `json:"billet"`
}

// PersonDetail is the full person detail returned by the search-by-txid endpoint.
type PersonDetail struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Name2       string        `json:"name2"`
	TxID        string        `json:"txId"`
	Email       string        `json:"email"`
	Status      int           `json:"status"`
	Phone       string        `json:"phone"`
	BirthDate   string        `json:"birthDate"`
	CellPhone   string        `json:"cellPhone"`
	MainAddress PersonAddress `json:"mainAddress"`
	Titles      []PersonTitle `json:"titles"`
}

// DocumentationType represents a document type available in the ERP.
type DocumentationType struct {
	Active       bool        `json:"active"`
	Code         string      `json:"code"`
	ID           int         `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	DefaultText  string      `json:"defaultText"`
	UseType      UseTypeItem `json:"useType"`
	SynSuiteCode *string     `json:"synsuiteCode"`
}

// UseTypeItem is a labeled enum value.
type UseTypeItem struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// ElectronicSignatureDocument represents an electronic signature document.
type ElectronicSignatureDocument struct {
	SignatureID                int     `json:"signatureId"`
	Description                string  `json:"description"`
	ContractID                 int     `json:"contractId"`
	AssignmentID               *int    `json:"assignmentId"`
	IntegrationCode            string  `json:"integrationCode"`
	PersonID                   int     `json:"personId"`
	PersonName                 string  `json:"personName"`
	IntegrationType            int     `json:"integrationType"`
	IntegrationTypeDescription string  `json:"integrationTypeDescription"`
	Status                     string  `json:"status"`
	Created                    string  `json:"created"`
	ExpirationDate             *string `json:"expirationDate"`
	LinkToSign                 string  `json:"linkToSign"`
}

// ResendSignatureRequest is the request body for resending an electronic signature.
// HTTP clients.
type ResendSignatureRequest struct {
	AssignmentID int `json:"assignmentId"`
	SignatureID  int `json:"signatureId"`
	DispatchType int `json:"dispatchType"`
}

// HolidayCheckRequest is the request body for validating if a date is a holiday.
type HolidayCheckRequest struct {
	Date string `json:"date"` // Format: YYYY-MM-DD
	UF   string `json:"UF"`
	City string `json:"city"`
}

// HolidayCheckResult is the response payload for the holiday check.
type HolidayCheckResult struct {
	IsHoliday bool   `json:"isHoliday"`
	Name      string `json:"name"`
	Date      string `json:"date"`
	City      string `json:"city"`
	Scope     string `json:"scope"`
	UF        string `json:"uf"`
}

// RegisterPerson creates a new person/customer record in the ERP.
//
// POST /external/integrations/thirdparty/people
func (c *Client) RegisterPerson(ctx context.Context, req RegisterPersonRequest) (*APIResponse[*RegisterPersonResult], error) {
	var result APIResponse[*RegisterPersonResult]
	err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/people"),
		req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePerson updates an existing person's contact information (email and/or phone).
//
// PUT /external/integrations/thirdparty/people/{personID}
func (c *Client) UpdatePerson(ctx context.Context, personID int, req UpdatePersonRequest) (*APIResponse[UpdatePersonResult], error) {
	var result APIResponse[UpdatePersonResult]
	err := c.doJSON(ctx, http.MethodPut,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/people/%d", personID)),
		req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCustomerAddress updates the main address of a customer.
//
// PUT /external/integrations/thirdparty/updateaddress/{clientID}
func (c *Client) UpdateCustomerAddress(ctx context.Context, clientID int, req UpdateAddressRequest) (*APIResponse[UpdateAddressResult], error) {
	var result APIResponse[UpdateAddressResult]
	err := c.doJSON(ctx, http.MethodPut,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/updateaddress/%d", clientID)),
		req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListCustomers returns a paginated list of customers.
//
// GET /external/integrations/thirdparty/getclient
func (c *Client) ListCustomers(ctx context.Context, params PaginationParams) (*APIResponse[PagedData[Customer]], error) {
	u := buildURL(c.apiURL("/external/integrations/thirdparty/getclient"), map[string]string{
		"page":     strconv.Itoa(params.Page),
		"pageSize": strconv.Itoa(params.PageSize),
	})

	var result APIResponse[PagedData[Customer]]
	if err := c.doJSON(ctx, http.MethodGet, u, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListLocations returns a paginated list of company locations/branches.
//
// GET /external/integrations/thirdparty/companiesplacespaged
func (c *Client) ListLocations(ctx context.Context, params PaginationParams) (*APIResponse[PagedData[Location]], error) {
	u := buildURL(c.apiURL("/external/integrations/thirdparty/companiesplacespaged"), map[string]string{
		"Page":     strconv.Itoa(params.Page),
		"PageSize": strconv.Itoa(params.PageSize),
		"Filter":   params.Filter,
		"OrderBy":  params.OrderBy,
	})

	var result APIResponse[PagedData[Location]]
	if err := c.doJSON(ctx, http.MethodGet, u, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPersonByTxID searches for a person by CPF or CNPJ (txId).
//
// GET /external/integrations/thirdparty/people/txid/{txId}
func (c *Client) GetPersonByTxID(ctx context.Context, txID string) (*APIResponse[PersonDetail], error) {
	u := c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/people/txid/%s", txID))

	var result APIResponse[PersonDetail]
	if err := c.doJSON(ctx, http.MethodGet, u, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDocumentationTypes returns available document classification types.
//
// GET /external/integrations/thirdparty/getdocumentationtypes
func (c *Client) ListDocumentationTypes(ctx context.Context, filter, orderBy string) (*APIResponse[PagedData[DocumentationType]], error) {
	u := buildURL(c.apiURL("/external/integrations/thirdparty/getdocumentationtypes"), map[string]string{
		"Filter":  filter,
		"OrderBy": orderBy,
	})

	var result APIResponse[PagedData[DocumentationType]]
	if err := c.doJSON(ctx, http.MethodGet, u, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetElectronicSignatureByContract retrieves electronic signature documents for a contract.
//
// GET /external/integrations/thirdparty/suite/electronicsignatures/getdocumentlinktosignbycontract/{contractID}
func (c *Client) GetElectronicSignatureByContract(ctx context.Context, contractID int) ([]ElectronicSignatureDocument, error) {
	u := c.apiURL(fmt.Sprintf(
		"/external/integrations/thirdparty/suite/electronicsignatures/getdocumentlinktosignbycontract/%d",
		contractID,
	))

	var result struct {
		Success  bool                          `json:"success"`
		Response []ElectronicSignatureDocument `json:"response"`
	}
	if err := c.doJSON(ctx, http.MethodGet, u, nil, &result); err != nil {
		return nil, err
	}
	return result.Response, nil
}

// ResendElectronicSignature resends an electronic signature document notification.
//
// POST /external/integrations/thirdparty/suite/electronicsignatures/resend
func (c *Client) ResendElectronicSignature(ctx context.Context, req ResendSignatureRequest) error {
	return c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/suite/electronicsignatures/resend"),
		req, nil)
}

// IsHoliday checks whether a given date is a holiday for a specific state/city.
//
// POST /external/integrations/thirdparty/isholiday
func (c *Client) IsHoliday(ctx context.Context, req HolidayCheckRequest) (*APIResponse[HolidayCheckResult], error) {
	var result APIResponse[HolidayCheckResult]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/isholiday"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
