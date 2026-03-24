// Package elleven provides a Go SDK for the ERP Elleven Third-Party API.
//
// The SDK covers all documented modules: authentication, Suite (people/customers),
// CRM, Service Desk, ISP/Telecom, Billing, Finance and Third-Party Billing.
//
// Basic usage:
//
//	client, err := NewClient(Config{
//	    BaseURL: "https://erp.example.com",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Authenticate
//	tokenResp, err := client.AuthenticateClientCredentials(ctx, ClientCredentialsRequest{
//	    ClientID:     "your-client-id",
//	    ClientSecret: "your-client-secret",
//	    SynData:      "your-syndata-token",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	client.SetTokenResponse(tokenResp)
package elleven

// APIResponse is the standard response envelope returned by all ERP Elleven API endpoints.
// T is the type of the response payload.
type APIResponse[T any] struct {
	Success          bool         `json:"success"`
	Messages         []APIMessage `json:"messages"`
	Response         T            `json:"response"`
	DataResponseType *string      `json:"dataResponseType"`
	ElapsedTime      *float64     `json:"elapsedTime"`
}

// APIMessage holds a message returned inside an API response.
type APIMessage struct {
	Message string `json:"message"`
	Code    *int   `json:"code"`
	Type    string `json:"type"` // "Success", "Error", "Warning"
}

// PagedData wraps paginated list responses from the API.
type PagedData[T any] struct {
	Data         []T `json:"data"`
	TotalRecords int `json:"totalRecords"`
	Page         int `json:"page"`
	PageSize     int `json:"pageSize"`
	TotalPages   int `json:"totalPages"`
}

// PaginationParams holds common query parameters for paginated endpoints.
type PaginationParams struct {
	Page     int
	PageSize int
	Filter   string
	OrderBy  string
}

// SimpleResponse wraps responses that only return a success flag and messages.
type SimpleResponse struct {
	Success  bool         `json:"success"`
	Messages []APIMessage `json:"messages"`
}
