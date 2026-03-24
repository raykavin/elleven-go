// Package main demonstrates the complete usage of the elleven Go SDK.
//
// This example covers: authentication, customer management, CRM operations,
// service desk, ISP/Telecom, billing, and finance.
//
// To run this example:
//
//	ELLEVEN_URL=http://erp.example.com \
//	ELLEVEN_CLIENT_ID=your-client-id \
//	ELLEVEN_CLIENT_SECRET=your-client-secret \
//	ELLEVEN_SYNDATA=your-syndata-token \
//	go run examples/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	elleven "github.com/raykavin/elleven-go"
)

func main() {
	// 1. Create the client
	baseURL := getenv("ELLEVEN_URL", "http://erp.example.com")

	client, err := elleven.NewClient(elleven.Config{
		BaseURL:      baseURL,
		APIPort:      "45715", // default, can be omitted
		AuthPort:     "45700", // default, can be omitted
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryWaitMin: time.Second,
		RetryWaitMax: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}

	ctx := context.Background()

	// 2. Authenticate
	fmt.Println("Authentication")

	// Option A: Modern client_credentials flow (recommended)
	tokenResp, err := client.AuthenticateClientCredentials(ctx, elleven.ClientCredentialsRequest{
		ClientID:     getenv("ELLEVEN_CLIENT_ID", ""),
		ClientSecret: getenv("ELLEVEN_CLIENT_SECRET", ""),
		SynData:      getenv("ELLEVEN_SYNDATA", ""),
	})
	if err != nil {
		log.Printf("client_credentials auth failed: %v", err)

		// Option B: Legacy password flow (fallback)
		tokenResp, err = client.AuthenticateLegacy(ctx, elleven.LegacyAuthRequest{
			Username: getenv("ELLEVEN_USER", ""),
			Password: getenv("ELLEVEN_PASS", ""),
			SynData:  getenv("ELLEVEN_SYNDATA", ""),
		})
		if err != nil {
			log.Fatalf("legacy auth failed: %v", err)
		}
	}

	// Store the token, all subsequent calls will include it automatically
	client.SetTokenResponse(tokenResp)
	fmt.Printf("Authenticated! Token expires at: %s\n", client.TokenExpiresAt().Format(time.RFC3339))

	// 3. Refresh token before expiry (optional)

	if time.Until(client.TokenExpiresAt()) < 5*time.Minute {
		newToken, err := client.RefreshToken(ctx, elleven.RefreshTokenRequest{
			RefreshToken: client.GetRefreshToken(),
		})
		if err != nil {
			log.Printf("token refresh failed: %v", err)
		} else {
			client.SetTokenResponse(newToken)
			fmt.Println("Token refreshed successfully")
		}
	}

	// 4. Suite & Customer Management

	fmt.Println("\n Suite - Customers ")

	// Register a new person
	regResp, err := client.RegisterPerson(ctx, elleven.RegisterPersonRequest{
		TypeTxID:     "CPF",
		TxID:         "12345678900",
		Name:         "João da Silva",
		Email:        "joao@example.com",
		Client:       true,
		Situation:    1,
		StreetType:   "Rua",
		PostalCode:   "97010001",
		Street:       "das Flores",
		Number:       "123",
		Neighborhood: "Centro",
		City:         "Santa Maria",
		State:        "RS",
		CodeCountry:  "1058",
	})
	if err != nil {
		if elleven.IsNotFound(err) {
			fmt.Println("Person not found")
		} else {
			log.Printf("RegisterPerson error: %v", err)
		}
	} else {
		fmt.Printf("RegisterPerson: success=%v, messages=%v\n", regResp.Success, regResp.Messages)
	}

	// Search by CPF/CNPJ
	personResp, err := client.GetPersonByTxID(ctx, "12345678900")
	if err != nil {
		log.Printf("GetPersonByTxID error: %v", err)
	} else {
		fmt.Printf("Found person: ID=%d, Name=%s\n", personResp.Response.ID, personResp.Response.Name)
	}

	// List customers with pagination
	customers, err := client.ListCustomers(ctx, elleven.PaginationParams{Page: 1, PageSize: 20})
	if err != nil {
		log.Printf("ListCustomers error: %v", err)
	} else {
		fmt.Printf("Customers: %d total, showing %d\n",
			customers.Response.TotalRecords, len(customers.Response.Data))
	}

	// Check if a date is a holiday
	holiday, err := client.IsHoliday(ctx, elleven.HolidayCheckRequest{
		Date: "2024-05-13",
		UF:   "RS",
		City: "Santa Maria",
	})
	if err != nil {
		log.Printf("IsHoliday error: %v", err)
	} else {
		if holiday.Response.IsHoliday {
			fmt.Printf("Holiday found: %s (%s)\n", holiday.Response.Name, holiday.Response.Scope)
		} else {
			fmt.Println("Not a holiday")
		}
	}

	// 5. CRM
	fmt.Println("\n CRM ")

	// Check viability at an address
	viab, err := client.VerifyViability(ctx, elleven.VerifyViabilityRequest{
		FullAddress: elleven.ViabilityAddress{
			Address:      "Rua das Flores",
			Number:       "123",
			Neighborhood: "Centro",
			City:         "Santa Maria",
			State:        "RS",
			PostalCode:   "97010001",
		},
		Distance: "200",
		Lead: elleven.ViabilityPerson{
			ID:   "1",
			Name: "João da Silva",
			TxID: "12345678900",
		},
	})
	if err != nil {
		log.Printf("VerifyViability error: %v", err)
	} else {
		fmt.Printf("Viability: %v (CTOs: %d, Ports: %d)\n",
			viab.Response.Viability, viab.Response.CTOs, viab.Response.Ports)
	}

	// Get contract types available
	ctypes, err := client.GetContractTypesAndServices(ctx)
	if err != nil {
		log.Printf("GetContractTypesAndServices error: %v", err)
	} else {
		for _, ct := range ctypes.Response {
			fmt.Printf("  Contract type: [%s] %s - collection days: %v\n",
				ct.Code, ct.Title, ct.CollectionDays)
		}
	}

	// Start a new sale
	saleResp, err := client.StartSale(ctx, elleven.StartSaleRequest{
		TxID:                                "12345678900",
		CompanyPlaceTxID:                    "00.000.000/0001-00",
		ContractType:                        "PF - Composto",
		CRMCampaignCode:                     "CAMP001",
		CRMPriceListCode:                    "PL001",
		ContractPaymentFormCode:             "BOLETO",
		ContractFinancialCollectionTypeCode: "COBRANCA",
		ContractCollectionDay:               "10",
		ServiceProducts: []elleven.SaleServiceProduct{
			{Code: "FIBER_100", Quantity: 1, Amount: 99.90, PaymentFormCode: "BOLETO"},
		},
		Observations: "Venda via integração",
	})
	if err != nil {
		log.Printf("StartSale error: %v", err)
	} else {
		fmt.Printf("Sale started! Protocol ID: %d\n", saleResp.Response.ProtocolID)

		// Cancel the sale (for demo purposes)
		cancelResp, err := client.CancelSale(ctx, elleven.CancelSaleRequest{
			ProtocolID:  saleResp.Response.ProtocolID,
			Description: "Demo cancellation",
		})
		if err != nil {
			log.Printf("CancelSale error: %v", err)
		} else {
			fmt.Printf("Sale cancelled: %s\n", cancelResp.Response)
		}
	}

	// 6. Service Desk
	fmt.Println("\n Service Desk ")

	solResp, err := client.OpenSimpleSolicitation(ctx, elleven.OpenSimpleSolicitationRequest{
		Description:          "Cliente sem sinal de internet",
		ClientID:             123456,
		ContractID:           789,
		ContractServiceTagID: 1,
		Close:                false,
	})
	if err != nil {
		log.Printf("OpenSimpleSolicitation error: %v", err)
	} else {
		fmt.Printf("Solicitation opened: success=%v\n", solResp.Success)
	}

	// 7. ISP / Telecom
	fmt.Println("\n ISP / Telecom ")

	statusResp, err := client.GetAccessPointStatusByContract(ctx, 789)
	if err != nil {
		log.Printf("GetAccessPointStatusByContract error: %v", err)
	} else {
		fmt.Printf("Access point: title=%s, active=%v, inMaintenance=%v\n",
			statusResp.Response.Title, statusResp.Response.Active, statusResp.Response.InMaintenance)
	}

	// Register user traffic
	trafficResp, err := client.RegisterUserTraffic(ctx, elleven.RegisterUserTrafficRequest{
		User:          "client@provider.com",
		Date:          time.Now().UTC().Format(time.RFC3339),
		Download:      524288000,   // 500 MB
		Upload:        104857600,   // 100 MB
		TotalDownload: 10737418240, // 10 GB
		TotalUpload:   2147483648,  // 2 GB
	})
	if err != nil {
		log.Printf("RegisterUserTraffic error: %v", err)
	} else {
		fmt.Printf("Traffic registered for user: %s\n", trafficResp.Data.User)
	}

	// 8. Finance
	fmt.Println("\n Finance ")

	// Get open invoices by CPF
	invoices, err := client.GetOpenInvoicesByTxID(ctx, "12345678900")
	if err != nil {
		log.Printf("GetOpenInvoicesByTxID error: %v", err)
	} else {
		fmt.Printf("Open invoices: %d found\n", len(invoices.Response))
		for _, inv := range invoices.Response {
			fmt.Printf("  Invoice #%d - Title: %s - Expires: %s - R$ %.2f\n",
				inv.ID, inv.Billet.Title, inv.Billet.ExpirationDate, inv.Billet.Amount.FinalValue)
		}
	}

	// Get all invoices by contract
	contractBillets, err := client.GetOpenInvoicesByContract(ctx, 789)
	if err != nil {
		log.Printf("GetOpenInvoicesByContract error: %v", err)
	} else {
		for _, b := range contractBillets.Response {
			fmt.Printf("  Contract billet: %s - expires: %s\n", b.Title, b.ExpirationDate)
		}
	}

	// Generate PIX for an invoice
	if invoices != nil && len(invoices.Response) > 0 {
		pixResp, err := client.GeneratePix(ctx, invoices.Response[0].ID)
		if err != nil {
			log.Printf("GeneratePix error: %v", err)
		} else {
			fmt.Printf("PIX generated: QR Code length=%d, total=R$%.2f\n",
				len(pixResp.Response.QRCode), pixResp.Response.TotalAmount)
		}
	}

	// 9. Error handling
	fmt.Println("\n Error Handling ")

	_, err = client.GetPersonByTxID(ctx, "00000000000") // non-existent
	if err != nil {
		switch {
		case elleven.IsNotFound(err):
			fmt.Println("Person not found (404)")
		case elleven.IsUnauthorized(err):
			fmt.Println("Authentication required (401) - call Authenticate* again")
		case elleven.IsRateLimited(err):
			fmt.Println("Rate limited (429) - wait and retry")
		case elleven.IsServerError(err):
			fmt.Println("Server error (5xx) - temporary issue")
		default:
			if apiErr, ok := err.(*elleven.APIError); ok {
				fmt.Printf("API error %d: %s\n", apiErr.StatusCode, apiErr.Error())
			} else {
				fmt.Printf("Transport/network error: %v\n", err)
			}
		}
	}

	fmt.Println("\nDone!")
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
