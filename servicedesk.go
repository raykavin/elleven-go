package elleven

import (
	"context"
	"fmt"
	"net/http"
)

// Service Desk

// AssignmentReport holds a report entry within a solicitation assignment.
type AssignmentReport struct {
	BeginningDate string `json:"beginningDate"` // RFC3339
	FinalDate     string `json:"finalDate"`     // RFC3339
	Description   string `json:"description"`
}

// SolicitationAssignment holds the task/assignment details for a solicitation.
type SolicitationAssignment struct {
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Priority       int              `json:"priority"`
	BeginningDate  string           `json:"beginningDate"` // RFC3339
	FinalDate      string           `json:"finalDate"`     // RFC3339
	Report         AssignmentReport `json:"report"`
	CompanyPlaceID int              `json:"companyPlaceId"`
}

// MatrixType distinguishes between external (1) and internal (2) solicitations.
type MatrixType int

const (
	MatrixTypeExternal MatrixType = 1
	MatrixTypeInternal MatrixType = 2
)

// OpenDetailedSolicitationRequest is the request body for creating a detailed solicitation.
//
// Use MatrixType = MatrixTypeExternal (1) for client-facing solicitations.
// Use MatrixType = MatrixTypeInternal (2) for internal solicitations (also requires
// TeamCode and AuthenticationAccessPointCode fields via the extended request).
type OpenDetailedSolicitationRequest struct {
	IncidentStatusID             int                    `json:"incidentStatusId"`
	PersonID                     int                    `json:"personId"`
	ClientID                     int                    `json:"clientId"`
	IncidentTypeID               int                    `json:"incidentTypeId"`
	ContractServiceTagID         int                    `json:"contractServiceTagId"`
	CatalogServiceID             int                    `json:"catalogServiceId"`
	ServiceLevelAgreementID      int                    `json:"serviceLevelAgreementId"`
	CatalogServiceItemID         int                    `json:"catalogServiceItemId"`
	CatalogServiceItemClassID    int                    `json:"catalogServiceItemClassId"`
	MatrixType                   MatrixType             `json:"matrixType"`
	Assignment                   SolicitationAssignment `json:"assignment"`
	ContractServiceTagCategory   string                 `json:"contractServiceTagCategory"`
	SolicitationServiceCategory1 string                 `json:"solicitationServiceCategory1"`
	SolicitationServiceCategory2 string                 `json:"solicitationServiceCategory2"`
	SolicitationServiceCategory3 string                 `json:"solicitationServiceCategory3"`
	SolicitationServiceCategory4 string                 `json:"solicitationServiceCategory4"`
	SolicitationServiceCategory5 string                 `json:"solicitationServiceCategory5"`

	// Internal solicitation fields (MatrixTypeInternal only)
	TeamCode                      string `json:"teamCode,omitempty"`
	AuthenticationAccessPointCode string `json:"authenticationAccessPointCode,omitempty"`
}

// OpenSimpleSolicitationRequest is the request body for opening a simple solicitation.
type OpenSimpleSolicitationRequest struct {
	Description          string `json:"description"`
	ClientID             int    `json:"clientId"`
	ContractID           int    `json:"contractId"`
	ContractServiceTagID int    `json:"contractServiceTagId"`
	Close                bool   `json:"close"`
}

// OpenCRMSolicitationRequest is the request body for opening a CRM-linked solicitation.
type OpenCRMSolicitationRequest struct {
	Description          string `json:"description"`
	ContractID           int    `json:"contractId"`
	ContractServiceTagID int    `json:"contractServiceTagId"`
}

// MaintainSolicitationRequest is the request body for maintaining (updating) a solicitation.
type MaintainSolicitationRequest struct {
	Protocol                      string `json:"protocol"`
	SolicitationServiceMatrixCode string `json:"solicitationServiceMatrixCode"`
	IncidentTypeCode              string `json:"incidentTypeCode"`
	Report                        string `json:"report"`
}

// AddNoteToSolicitationRequest is the request body for adding a note to a solicitation.
type AddNoteToSolicitationRequest struct {
	AssignmentID int    `json:"assignmentId"`
	PersonID     int    `json:"personId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

// CreateSolicitationReportRequest is the request body for creating a report on a solicitation.
type CreateSolicitationReportRequest struct {
	AssignmentID       int    `json:"assignmentId"`
	Protocol           int    `json:"protocol"`
	IncidentStatusID   int    `json:"incidentStatusId"`
	Description        string `json:"description"`
	Progress           int    `json:"progress"`
	Priority           int    `json:"priority"`
	NotificationTarget int    `json:"notificationTarget"`
	PrivateReport      bool   `json:"privateReport"`
}

// OpenDetailedSolicitation creates a detailed (external or internal) solicitation.
//
// POST /external/integrations/thirdparty/opendetailedsolicitation
func (c *Client) OpenDetailedSolicitation(ctx context.Context, req OpenDetailedSolicitationRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/opendetailedsolicitation"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OpenSimpleSolicitation creates a basic solicitation linked to a client and contract.
//
// POST /external/integrations/thirdparty/opensolicitation
func (c *Client) OpenSimpleSolicitation(ctx context.Context, req OpenSimpleSolicitationRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/opensolicitation"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OpenCRMSolicitation creates a CRM-linked solicitation.
//
// POST /external/integrations/thirdparty/opensolicitationcrm
func (c *Client) OpenCRMSolicitation(ctx context.Context, req OpenCRMSolicitationRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/opensolicitationcrm"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OpenSolicitationForAccessPoint opens a solicitation for a specific access point
// identified by its integration code.
//
// POST /external/integrations/thirdparty/opensolicitationpopeventerror/{integrationCode}
func (c *Client) OpenSolicitationForAccessPoint(ctx context.Context, integrationCode string) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/opensolicitationpopeventerror/%s", integrationCode)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MaintainSolicitation updates an existing solicitation (status, type, report).
//
// POST /external/integrations/thirdparty/solicitationmaintenance
func (c *Client) MaintainSolicitation(ctx context.Context, req MaintainSolicitationRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/solicitationmaintenance"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CloseSolicitationForAccessPoint closes a solicitation for a specific access point.
//
// POST /external/integrations/thirdparty/closesolicitationpopeventerror/{integrationCode}
func (c *Client) CloseSolicitationForAccessPoint(ctx context.Context, integrationCode string) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/closesolicitationpopeventerror/%s", integrationCode)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OpenSolicitationForClient opens a solicitation triggered by a client error event.
//
// POST /external/integrations/thirdparty/opensolicitationclienteventerror/{clientID}
func (c *Client) OpenSolicitationForClient(ctx context.Context, clientID int) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/opensolicitationclienteventerror/%d", clientID)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CloseSolicitationForClient closes a solicitation triggered by a client error event.
//
// POST /external/integrations/thirdparty/closesolicitationclienteventerror/{clientID}
func (c *Client) CloseSolicitationForClient(ctx context.Context, clientID int) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/closesolicitationclienteventerror/%d", clientID)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddNoteToSolicitation adds a note/comment to an existing solicitation.
//
// POST /external/integrations/thirdparty/solicitationnewnote
func (c *Client) AddNoteToSolicitation(ctx context.Context, req AddNoteToSolicitationRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/solicitationnewnote"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateSolicitationReport creates a progress report on a solicitation and optionally
// updates its status.
//
// POST /external/integrations/thirdparty/projects/createsolicitationreport
func (c *Client) CreateSolicitationReport(ctx context.Context, req CreateSolicitationReportRequest) (*APIResponse[any], error) {
	var result APIResponse[any]
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/projects/createsolicitationreport"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadAttachmentToSolicitation uploads a file attachment to a solicitation.
//
// POST /external/integrations/thirdparty/projects/assignmentsuploads
// Query params: protocol, documentationTypeCode
func (c *Client) UploadAttachmentToSolicitation(
	ctx context.Context,
	protocol int,
	documentationTypeCode string,
	fileName string,
	fileData []byte,
) (*APIResponse[any], error) {
	u := buildURL(
		c.apiURL("/external/integrations/thirdparty/projects/assignmentsuploads"),
		map[string]string{
			"protocol":              fmt.Sprintf("%d", protocol),
			"documentationTypeCode": documentationTypeCode,
		},
	)

	var result APIResponse[any]
	if err := c.doMultipart(ctx, u, "File", fileName, fileData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadAttachmentFromSolicitation downloads a file attachment from a solicitation.
//
// Returns the raw file bytes and the Content-Type header value.
//
// GET /external/integrations/thirdparty/projects/assignmentsuploads/getfile
func (c *Client) DownloadAttachmentFromSolicitation(ctx context.Context, attachmentID, assignmentID int) ([]byte, string, error) {
	u := buildURL(
		c.apiURL("/external/integrations/thirdparty/projects/assignmentsuploads/getfile"),
		map[string]string{
			"id":           fmt.Sprintf("%d", attachmentID),
			"assignmentId": fmt.Sprintf("%d", assignmentID),
		},
	)
	return c.doDownload(ctx, u)
}
