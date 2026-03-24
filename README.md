# Unofficial Elleven SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/raykavin/ellevensdk.svg)](https://pkg.go.dev/github.com/raykavin/ellevensdk)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/raykavin/ellevensdk)](https://goreportcard.com/report/github.com/raykavin/ellevensdk)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go SDK for the **Elleven ERP Third-Party API**.

Covers all documented modules: Authentication, Suite (People/Customers), CRM, Service Desk, ISP/Telecom, Billing, Finance, and Third-Party Billing.

---

## Installation

```bash
go get github.com/raykavin/ellevensdk
```

Requires **Go 1.22+**. No external dependencies, uses only the Go standard library.

---

## Quick Start

```go
import elleven "github.com/raykavin/ellevensdk"

// 1. Create the client
client, err := elleven.NewClient(elleven.Config{
    BaseURL: "http://erp.yourcompany.com",
})
if err != nil {
    log.Fatal(err)
}

// 2. Authenticate (client_credentials flow)
ctx := context.Background()
tokenResp, err := client.AuthenticateClientCredentials(ctx, elleven.ClientCredentialsRequest{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    SynData:      "your-syndata-token",
})
if err != nil {
    log.Fatal(err)
}
client.SetTokenResponse(tokenResp)

// 3. Use any resource
customers, err := client.ListCustomers(ctx, elleven.PaginationParams{Page: 1, PageSize: 20})
```

---

## Configuration

```go
client, err := elleven.NewClient(elleven.Config{
    // Required
    BaseURL: "http://erp.yourcompany.com",

    // Optional defaults shown
    APIPort:      "45715",           // port for all API calls
    AuthPort:     "45700",           // port for authentication
    Timeout:      30 * time.Second,  // per-request timeout
    MaxRetries:   3,                 // retries on 429/502/503/504
    RetryWaitMin: time.Second,       // min wait between retries
    RetryWaitMax: 30 * time.Second,  // max wait between retries (exponential cap)
    AccessToken:  "",                // pre-configure a token directly

    // Optional inject custom *http.Client (for TLS, proxies, etc.)
    HTTPClient: &http.Client{Transport: myTransport},
})
```

---

## Authentication

The ERP Elleven API uses OAuth2. Three flows are supported:

### Modern client_credentials (recommended)

Credentials come from **Settings / Users / Integration User** in the ERP interface.

```go
tokenResp, err := client.AuthenticateClientCredentials(ctx, elleven.ClientCredentialsRequest{
    ClientID:     "client-id-from-erp",
    ClientSecret: "client-secret-from-erp",
    SynData:      "syndata-token",
})
client.SetTokenResponse(tokenResp)
```

### Legacy password grant

Uses fixed platform credentials (`synauth`). Requires the integration user's
username and password.

```go
tokenResp, err := client.AuthenticateLegacy(ctx, elleven.LegacyAuthRequest{
    Username: "integration-user",
    Password: "password",
    SynData:  "syndata-token",
})
client.SetTokenResponse(tokenResp)
```

### Refresh Token

```go
newToken, err := client.RefreshToken(ctx, elleven.RefreshTokenRequest{
    RefreshToken: client.GetRefreshToken(),
})
client.SetTokenResponse(newToken)
```

### Token lifecycle

```go
// Check expiry before making calls
if time.Until(client.TokenExpiresAt()) < 5*time.Minute {
    // refresh proactively
}
```

---

## Modules

### Suite People / Customers

```go
// Register a new person
resp, err := client.RegisterPerson(ctx, elleven.RegisterPersonRequest{
    TypeTxID: "CPF", TxID: "12345678900",
    Name: "João Silva", Email: "joao@example.com",
    Client: true, Situation: 1,
    PostalCode: "97010001", Street: "das Flores",
    Number: "123", City: "Santa Maria", State: "RS",
})

// Update contact info
resp, err := client.UpdatePerson(ctx, personID, elleven.UpdatePersonRequest{
    ID: personID,
    Email: elleven.UpdatePersonEmailRequest{Email: "new@email.com"},
    Phone: elleven.UpdatePersonPhoneRequest{CellPhone: "55999990000", CellPhoneHasWhatsapp: true},
})

// Update address
resp, err := client.UpdateCustomerAddress(ctx, clientID, elleven.UpdateAddressRequest{
    ZipCode: "97010001", PublicPlace: "Rua das Flores",
    Number: "456", City: "Santa Maria", State: "RS",
})

// List customers (paginated)
page, err := client.ListCustomers(ctx, elleven.PaginationParams{Page: 1, PageSize: 50})
fmt.Println(page.Response.TotalRecords, "total customers")

// List company locations
locs, err := client.ListLocations(ctx, elleven.PaginationParams{
    Page: 1, PageSize: 10, Filter: "matriz",
})

// Search by CPF/CNPJ
person, err := client.GetPersonByTxID(ctx, "12345678900")

// List documentation types
dtypes, err := client.ListDocumentationTypes(ctx, "", "title")

// Electronic signatures
sigs, err := client.GetElectronicSignatureByContract(ctx, contractID)
err = client.ResendElectronicSignature(ctx, elleven.ResendSignatureRequest{
    SignatureID: sigs[0].SignatureID, DispatchType: 1,
})

// Holiday check
h, err := client.IsHoliday(ctx, elleven.HolidayCheckRequest{
    Date: "2024-12-25", UF: "RS", City: "Santa Maria",
})
```

---

### CRM

```go
// Create a lead
resp, err := client.CreateLead(ctx, elleven.CreateLeadRequest{
    PersonalData: elleven.LeadPersonalData{Name: "Maria", TxID: "98765432100"},
    Contact:      elleven.LeadContact{CellPhone: "55999990000", Email: "maria@example.com"},
    Address:      elleven.LeadAddress{PostalCode: "97010001", City: "Santa Maria"},
})

// Check service coverage viability
v, err := client.VerifyViability(ctx, elleven.VerifyViabilityRequest{
    FullAddress: elleven.ViabilityAddress{
        Address: "Rua das Flores", Number: "123",
        City: "Santa Maria", State: "RS", PostalCode: "97010001",
    },
    Distance: "200",
    Lead: elleven.ViabilityPerson{Name: "João", TxID: "12345678900"},
})
fmt.Println("Has coverage:", v.Response.Viability)

// Get available contract types
ctypes, err := client.GetContractTypesAndServices(ctx)

// Get campaigns and price lists
campaigns, err := client.GetCampaignsAndPriceListServices(ctx)

// Start a new negotiation/sale
sale, err := client.StartSale(ctx, elleven.StartSaleRequest{
    TxID:             "12345678900",
    CompanyPlaceTxID: "00.000.000/0001-00",
    ContractType:     "PF - Composto",
    ServiceProducts: []elleven.SaleServiceProduct{
        {Code: "FIBER_100", Quantity: 1, Amount: 99.90},
    },
})
fmt.Println("Protocol:", sale.Response.ProtocolID)

// Cancel a negotiation
_, err = client.CancelSale(ctx, elleven.CancelSaleRequest{
    ProtocolID: sale.Response.ProtocolID, Description: "Customer gave up",
})

// Add service to contract
_, err = client.AddContractServices(ctx, elleven.AddContractServicesRequest{
    ContractNumber: "12345",
    NewServices: []elleven.NewContractService{
        {Code: "SVC_TV", Quantity: 1, Price: 29.90},
    },
})

// Remove service from contract
_, err = client.RemoveContractServices(ctx, elleven.RemoveContractServicesRequest{
    ContractNumber: "12345",
    RemoveService:  []elleven.RemoveContractItem{{ContractItem: 101}},
})
```

---

### Service Desk

```go
// Open a simple solicitation
resp, err := client.OpenSimpleSolicitation(ctx, elleven.OpenSimpleSolicitationRequest{
    Description: "No internet signal", ClientID: 123, ContractID: 456,
    ContractServiceTagID: 1,
})

// Open a detailed external solicitation
resp, err := client.OpenDetailedSolicitation(ctx, elleven.OpenDetailedSolicitationRequest{
    ClientID: 123, IncidentTypeID: 5,
    MatrixType: elleven.MatrixTypeExternal,
    Assignment: elleven.SolicitationAssignment{
        Title: "Network down", Priority: 1,
    },
})

// Maintain (update) a solicitation
_, err = client.MaintainSolicitation(ctx, elleven.MaintainSolicitationRequest{
    Protocol: "PROT-001", Report: "Technician dispatched",
})

// Add a note
_, err = client.AddNoteToSolicitation(ctx, elleven.AddNoteToSolicitationRequest{
    AssignmentID: 789, Title: "Update", Description: "Signal restored",
})

// Create progress report
_, err = client.CreateSolicitationReport(ctx, elleven.CreateSolicitationReportRequest{
    AssignmentID: 789, Protocol: 123, Progress: 75, Description: "Working on it",
})

// Upload attachment
_, err = client.UploadAttachmentToSolicitation(ctx, 123, "1.01", "photo.jpg", imageBytes)

// Download attachment
data, contentType, err := client.DownloadAttachmentFromSolicitation(ctx, attachmentID, assignmentID)

// Event-based operations (access point / client)
_, err = client.OpenSolicitationForAccessPoint(ctx, "integration-code-abc")
_, err = client.CloseSolicitationForAccessPoint(ctx, "integration-code-abc")
_, err = client.OpenSolicitationForClient(ctx, clientID)
_, err = client.CloseSolicitationForClient(ctx, clientID)
```

---

### ISP / Telecom

```go
// Update connection technical data
_, err := client.UpdateConnection(ctx, connectionID, elleven.UpdateConnectionRequest{
    Mac: "AA:BB:CC:DD:EE:FF", TechnologyType: 1,
    OltID: 5, SlotOlt: 1, PortOlt: 10,
})

// Access point status
status, err := client.GetAccessPointStatusByContract(ctx, contractID)
fmt.Println("Active:", status.Response.Active)

statuses, err := client.GetAccessPointStatusByClient(ctx, clientID)
for _, s := range statuses.Response {
    fmt.Printf("%s: active=%v\n", s.Title, s.Active)
}

// Register user traffic data
_, err = client.RegisterUserTraffic(ctx, elleven.RegisterUserTrafficRequest{
    User:          "user@provider",
    Date:          time.Now().UTC().Format(time.RFC3339),
    Download:      524288000,
    Upload:        104857600,
    TotalDownload: 10737418240,
    TotalUpload:   2147483648,
})
```

---

### Billing (Faturamento)

```go
// Create contract via billing module
contract, err := client.CreateContract(ctx, elleven.CreateContractRequest{
    ClientInformation:   elleven.ContractClientInfo{TxID: "12345678900"},
    ContractInformation: elleven.ContractInfo{
        ContractType: "PF", CompanyPlaceTxID: "00.000.000/0001-00",
        CollectionDay: 10, PaymentFormCode: "BOLETO",
    },
    ServicesInformation: []elleven.ContractServiceInfo{
        {Code: "FIBER_100", Quantity: 1, Price: 99.90},
    },
})
fmt.Println("Contract:", contract.Response.ContractNumber)

// Approve contract
approval, err := client.ApproveContract(ctx, contractID, elleven.ApproveContractRequest{
    ProportionalityType:           0,
    ChangeBeginningDateOnApproval: true,
    ApprovalDate:                  time.Now().Format(time.RFC3339),
})

// Generate sale order
_, err = client.GenerateSaleOrder(ctx, elleven.GenerateSaleOrderRequest{
    ClientTxID: "12345678900", ContractNumber: "12345",
    Items: []elleven.SaleOrderItem{{ID: "product-uuid", Quantity: 1, UnitValue: 99.90}},
})

// Create eventual (one-time) billing value
_, err = client.CreateEventualValue(ctx, elleven.CreateEventualValueRequest{
    Type: elleven.EventualValueTypeDebit, ContractID: 12345,
    Description: "Late fee", UnitAmount: 15.00, Units: 1,
    MonthYear: "2024-12-01",
})

// Document management
docs, err := client.ListDocumentsByContract(ctx, "12345", documentationTypeID)
_, err = client.UploadAttachmentToContract(ctx, "12345", "1.01", "contract.pdf", pdfBytes)
data, ct, err := client.DownloadContractAttachment(ctx, attachmentID, contractID)

// Unlock a suspended contract
_, err = client.UnlockContract(ctx, "12345")

// Download DANFE (NF-e PDF)
pdfBytes, _, err := client.DownloadDANFE(ctx, documentID)
```

---

### Finance (Financeiro)

```go
// Get open invoices by CPF/CNPJ
invoices, err := client.GetOpenInvoicesByTxID(ctx, "12345678900")

// Get all invoices (any status) by CPF/CNPJ
all, err := client.GetAllInvoicesByTxID(ctx, "12345678900")

// Get open invoices by contract
billets, err := client.GetOpenInvoicesByContract(ctx, contractID)

// Download invoice PDF (boleto)
pdfBytes, _, err := client.GetInvoicePDF(ctx, invoiceID)

// Register a payment (settle an invoice)
payment, err := client.RegisterPayment(ctx, elleven.RegisterPaymentRequest{
    TransactionID:              "unique-uuid-v4",
    FinancialReceivableTitleID: invoiceID,
    PaidAmount:                 99.90,
    PaymentFormCode:            "PIX",
    ReceiptDate:                "2024-05-01",
    ClientPaidDate:             "2024-05-01",
})
fmt.Println("SynGW Transaction:", payment.Response.SynGwTransactionID)

// Consult payment processing status
status, err := client.ConsultPaymentStatus(ctx, synGWTransactionID)
fmt.Println("Status:", status.Response.Status.Label)

// Generate PIX QR Code for an invoice
pix, err := client.GeneratePix(ctx, invoiceID)
fmt.Println("PIX QR Code:", pix.Response.QRCode)
fmt.Printf("Total: R$ %.2f\n", pix.Response.TotalAmount)

// Renegotiation
info, err := client.GetRenegotiationInfo(ctx, elleven.RenegotiationInfoRequest{
    ReceivableIDs: []int{1, 2, 3}, Date: "2024-12-31",
})

result, err := client.RegisterRenegotiation(ctx, elleven.RegisterRenegotiationRequest{
    ReceivableIDs: []int{1, 2, 3},
    Parcels: []elleven.RenegotiationParcel{
        {Number: 1, ExpirationDate: "2024-12-01"},
        {Number: 2, ExpirationDate: "2025-01-01"},
    },
})
```

---

### Third-Party Billing (Co-faturamento)

```go
_, err := client.ConfirmInvoiceIssuance(ctx, elleven.ConfirmInvoiceIssuanceRequest{
    ID:        "internal-record-id",
    AccessKey: "NFe-access-key-44-chars",
    IssueDate: "2024-05-15",
})
```

---

## Error Handling

All methods return a `*APIError` on HTTP-level failures.

```go
_, err := client.GetPersonByTxID(ctx, "12345678900")
if err != nil {
    switch {
    case elleven.IsNotFound(err):
        // HTTP 404
    case elleven.IsUnauthorized(err):
        // HTTP 401 re-authenticate
    case elleven.IsForbidden(err):
        // HTTP 403 insufficient permissions
    case elleven.IsRateLimited(err):
        // HTTP 429 back off and retry
    case elleven.IsServerError(err):
        // HTTP 5xx transient server error
    default:
        if apiErr, ok := err.(*elleven.APIError); ok {
            fmt.Printf("Status: %d\n", apiErr.StatusCode)
            fmt.Printf("Messages: %v\n", apiErr.Messages)
            fmt.Printf("Raw body: %s\n", apiErr.RawBody)
        }
    }
}
```

---

## Pagination

Paginated endpoints accept `PaginationParams`:

```go
for page := 1; ; page++ {
    result, err := client.ListCustomers(ctx, elleven.PaginationParams{
        Page: page, PageSize: 100,
    })
    if err != nil {
        break
    }
    // process result.Response.Data ...
    if page >= result.Response.TotalPages {
        break
    }
}
```

---

## Retry

The client automatically retries on HTTP `429`, `502`, `503`, `504` with
exponential backoff. Configure via `Config.MaxRetries`, `RetryWaitMin`,
`RetryWaitMax`.

```go
client, _ := elleven.NewClient(elleven.Config{
    BaseURL:      "...",
    MaxRetries:   5,
    RetryWaitMin: 500 * time.Millisecond,
    RetryWaitMax: 60 * time.Second,
})
```

Set `MaxRetries: -1` disables retries (only the initial attempt is made).
The default is `MaxRetries: 3`.

---

## File Upload / Download

```go
// Upload
fileBytes, _ := os.ReadFile("document.pdf")
_, err := client.UploadAttachmentToContract(ctx, "12345", "1.01", "document.pdf", fileBytes)

// Download
data, contentType, err := client.DownloadContractAttachment(ctx, attachmentID, contractID)
if err == nil {
    os.WriteFile("output.pdf", data, 0644)
    fmt.Println("Content-Type:", contentType)
}
```

---

## Integration User Setup

1. In the ERP, go to **Settings / Users**.
2. Create a dedicated user and check the **"Usuário integrador"** option.
3. Each integration system should use its own integration user for auditing.
4. Obtain the `syndata` token from **Suite / Settings / Parameters / Integration/Map**.
5. For client_credentials, retrieve `Client Id` and `Client Secret` from the user record.

> **Note:** Integration users cannot log in to the ERP interface - they are API-only.

---

## Running Tests

```bash
cd ellevensdk
go test ./... -v
```

For integration tests against a real server:

```bash
ELLEVEN_URL=http://erp.example.com \
ELLEVEN_CLIENT_ID=... \
ELLEVEN_CLIENT_SECRET=... \
ELLEVEN_SYNDATA=... \
go test ./... -v -tags=integration
```

---

## Known Ambiguities / Assumptions

| # | Issue | Decision |
|---|-------|----------|
| 1 | `ResendElectronicSignature` is documented as `GET` with a JSON body | Implemented as `POST` - GET with body is non-standard and blocked by many HTTP clients |
| 2 | `GetRenegotiationInfo` is documented as `GET` with a JSON body | Implemented as `GET` with body to match documentation - some servers support this |
| 3 | `CreateContract` response uses `"Success"` (capital S) | Mapped faithfully with `json:"Success"` |
| 4 | `AddContractServices` response uses `"Sucess"` (typo) | Mapped faithfully with `json:"Sucess"` to match API |
| 5 | `Customer.id` field comes as string or number from the API | Typed as `any` to handle both |
| 6 | `StartSale` nested objects (`clientSaleInfo`, `activationInformation`, etc.) are not fully documented | Typed as `any` for flexibility |
| 7 | The API uses different URL variable cases (`{{URL}}` vs `{{url}}`) | Treated as the same variable |
| 8 | Base URL is tenant-specific (no fixed hostname) | `Config.BaseURL` must be provided by the caller |

---

## Endpoints Coverage

### Authentication (3/3)
- [x] `POST /connect/token` - Legacy (password grant)
- [x] `POST /connect/token` - Client credentials
- [x] `POST /connect/token` - Refresh token

### Suite (10/10)
- [x] `POST /external/integrations/thirdparty/people` - Register person
- [x] `PUT /external/integrations/thirdparty/people/{id}` - Update person
- [x] `PUT /external/integrations/thirdparty/updateaddress/{id}` - Update address
- [x] `GET /external/integrations/thirdparty/getclient` - List customers
- [x] `GET /external/integrations/thirdparty/companiesplacespaged` - List locations
- [x] `GET /external/integrations/thirdparty/people/txid/{txId}` - Search by CPF/CNPJ
- [x] `GET /external/integrations/thirdparty/getdocumentationtypes` - List doc types
- [x] `GET /external/integrations/thirdparty/suite/electronicsignatures/...` - Get e-signatures
- [x] `GET→POST /external/integrations/thirdparty/suite/electronicsignatures/resend` - Resend
- [x] `POST /external/integrations/thirdparty/isholiday` - Holiday check

### CRM (9/9)
- [x] `POST /external/crm/leads/create` - Create lead
- [x] `POST /external/integrations/thirdparty/crm/startsale` - Start sale
- [x] `POST /external/integrations/thirdparty/contract/addcontractservices`
- [x] `POST /external/integrations/thirdparty/contract/changecontractservices`
- [x] `POST /external/integrations/thirdparty/contract/removecontractitems`
- [x] `POST /external/integrations/thirdparty/verifyviability`
- [x] `POST /external/integrations/thirdparty/crm/cancelsale`
- [x] `GET /external/integrations/thirdparty/crm/contracttypesandservices`
- [x] `GET /external/integrations/thirdparty/crm/campaignsandpricelistservices`

### Service Desk (13/13)
- [x] `POST .../opendetailedsolicitation` - External (matrixType=1)
- [x] `POST .../opendetailedsolicitation` - Internal (matrixType=2)
- [x] `POST .../opensolicitation`
- [x] `POST .../opensolicitationcrm`
- [x] `POST .../opensolicitationpopeventerror/{code}`
- [x] `POST .../solicitationmaintenance`
- [x] `POST .../closesolicitationpopeventerror/{code}`
- [x] `POST .../opensolicitationclienteventerror/{id}`
- [x] `POST .../closesolicitationclienteventerror/{id}`
- [x] `POST .../solicitationnewnote`
- [x] `POST .../projects/createsolicitationreport`
- [x] `POST .../projects/assignmentsuploads` - Upload attachment
- [x] `GET .../projects/assignmentsuploads/getfile` - Download attachment

### ISP / Telecom (4/4)
- [x] `PUT /external/integrations/thirdparty/updateconnection/{id}`
- [x] `GET .../getaccesspointstatusbycontract/{id}`
- [x] `GET .../getaccesspointstatusbyclient/{id}`
- [x] `POST .../isp/authentication_contracts/traffic`

### Billing (9/9)
- [x] `POST /external/billing/contracts/create`
- [x] `POST .../salerequest`
- [x] `POST .../contract/contracteventualvalues`
- [x] `GET .../getfilesbycontract/{n}/documentationtype/{id}`
- [x] `POST .../contract/contractuploads` - Upload
- [x] `POST .../contracts/unlock/{n}`
- [x] `GET .../contract/contractuploads/getfile` - Download
- [x] `PUT .../contracts/approve/{id}`
- [x] `GET .../billings/invoices/download/{id}` - Download DANFE

### Finance (9/9)
- [x] `GET .../getopentitlesbytxid/{txId}`
- [x] `GET .../gettitlesbytxid/{txId}`
- [x] `GET .../getcontractbillets/{id}`
- [x] `GET .../GetBillet/{id}` - PDF download
- [x] `GET .../consultpayment`
- [x] `POST .../receivepayment`
- [x] `POST .../billings/registerpix`
- [x] `GET .../financial/getrenegotiationsinformations`
- [x] `POST .../financial/registerrenegotiations`

### Third-Party Billing (1/1)
- [x] `POST .../thirdparty_billing/invoice_note/confirm`

**Total: 58 endpoints covered**


---

## Contributing

Contributions to ellevensdk are welcome! Here are some ways you can help improve the project:

- **Report bugs and suggest features** by opening issues on GitHub
- **Submit pull requests** with bug fixes or new features
- **Improve documentation** to help other users and developers
- **Share your custom strategies** with the community

## License

ellevensdk is distributed under the **MIT License**.  
For complete license terms and conditions, see the [LICENSE](LICENSE) file in the repository.

---

## Contact

For support, collaboration, or questions about ellevensdk:

**Email**: [raykavin.meireles@gmail.com](mailto:raykavin.meireles@gmail.com)  
**GitHub**: [@raykavin](https://github.com/raykavin)  
