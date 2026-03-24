package elleven

import (
	"context"
	"net/http"
)

// Third-Party / Co-Billing (Faturamento Terceiro / Cofaturamento)

// ConfirmInvoiceIssuanceRequest is the request body for confirming that a
// third-party invoice (nota fiscal) was issued.
type ConfirmInvoiceIssuanceRequest struct {
	// ID is the internal identifier of the invoice issuance record.
	ID string `json:"id"`
	// AccessKey is the NF-e / NFS-e access key (chave de acesso).
	AccessKey string `json:"accessKey"`
	// IssueDate is the date the invoice was issued (format: YYYY-MM-DD).
	IssueDate string `json:"issueDate"`
}

// ConfirmInvoiceIssuance notifies the ERP that a third-party invoice was issued,
// providing the access key and issue date.
//
// POST /external/integrations/thirdparty/thirdparty_billing/invoice_note/confirm
func (c *Client) ConfirmInvoiceIssuance(ctx context.Context, req ConfirmInvoiceIssuanceRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/thirdparty_billing/invoice_note/confirm"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
