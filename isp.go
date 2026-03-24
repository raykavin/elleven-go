package elleven

import (
	"context"
	"fmt"
	"net/http"
)

// ISP / Telecom

// UpdateConnectionRequest is the request body for updating a connection's technical data.
type UpdateConnectionRequest struct {
	ID                          int    `json:"id"`
	FiberMac                    string `json:"fiberMac"`
	Mac                         string `json:"mac"`
	Password                    string `json:"password"`
	EquipmentType               int    `json:"equipmentType"`
	OltID                       int    `json:"oltId"`
	SlotOlt                     int    `json:"slotOlt"`
	PortOlt                     int    `json:"portOlt"`
	EquipmentSerialNumber       string `json:"equipmentSerialNumber"`
	IPType                      int    `json:"ipType"`
	EquipmentUser               string `json:"equipmentUser"`
	EquipmentPassword           string `json:"equipmentPassword"`
	AuthenticationSplitterID    int    `json:"authenticationSplitterId"`
	Port                        int    `json:"port"`
	WifiName                    string `json:"wifiName"`
	WifiPassword                string `json:"wifiPassword"`
	TechnologyType              int    `json:"technologyType"`
	AuthenticationAccessPointID int    `json:"authenticationAccessPointId"`
	UpdateConnectionParameter   bool   `json:"updateConnectionParameter"`
	ShouldMacUpdate             bool   `json:"shouldMacUpdate"`
	User                        string `json:"user"`
	IsIPoE                      bool   `json:"isIPoE"`
	Complement                  string `json:"complement"`
}

// AccessPointStatus represents the current status of an access point (e.g., OLT, POP).
type AccessPointStatus struct {
	Title         string `json:"title"`
	InMaintenance bool   `json:"inMaintenance"`
	Maintenance   any    `json:"maintenance"` // null or maintenance details
	Active        bool   `json:"active"`
}

// RegisterUserTrafficRequest is the request body for registering ISP user traffic data.
type RegisterUserTrafficRequest struct {
	// User is the PPPoE/IPoE username (e.g., "client@provider").
	User string `json:"user"`
	// Date is the timestamp for the traffic record (RFC3339).
	Date string `json:"date"`
	// Download is the current session download in bytes.
	Download int64 `json:"download"`
	// Upload is the current session upload in bytes.
	Upload int64 `json:"upload"`
	// TotalDownload is the total accumulated download in bytes.
	TotalDownload int64 `json:"totalDownload"`
	// TotalUpload is the total accumulated upload in bytes.
	TotalUpload int64 `json:"totalUpload"`
}

// RegisterUserTrafficResult is the response payload for traffic registration.
type RegisterUserTrafficResult struct {
	User          string `json:"user"`
	Date          string `json:"date"`
	Download      int64  `json:"download"`
	Upload        int64  `json:"upload"`
	TotalDownload int64  `json:"totalDownload"`
	TotalUpload   int64  `json:"totalUpload"`
}

// TrafficResponse wraps the traffic registration response.
type TrafficResponse struct {
	Success     bool                      `json:"success"`
	Data        RegisterUserTrafficResult `json:"data"`
	Messages    []APIMessage              `json:"messages"`
	ElapsedTime float64                   `json:"elapsedTime"`
}

// UpdateConnection updates the technical parameters of a connection.
//
// PUT /external/integrations/thirdparty/updateconnection/{connectionID}
func (c *Client) UpdateConnection(ctx context.Context, connectionID int, req UpdateConnectionRequest) (*APIResponse[bool], error) {
	var result APIResponse[bool]
	if err := c.doJSON(ctx, http.MethodPut,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/updateconnection/%d", connectionID)),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAccessPointStatusByContract returns the access point status for a given contract.
//
// GET /external/integrations/thirdparty/getaccesspointstatusbycontract/{contractID}
func (c *Client) GetAccessPointStatusByContract(ctx context.Context, contractID int) (*APIResponse[AccessPointStatus], error) {
	var result APIResponse[AccessPointStatus]
	if err := c.doJSON(ctx, http.MethodGet,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/getaccesspointstatusbycontract/%d", contractID)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAccessPointStatusByClient returns all access point statuses for a given client.
//
// GET /external/integrations/thirdparty/getaccesspointstatusbyclient/{clientID}
func (c *Client) GetAccessPointStatusByClient(ctx context.Context, clientID int) (*APIResponse[[]AccessPointStatus], error) {
	var result APIResponse[[]AccessPointStatus]
	if err := c.doJSON(ctx, http.MethodGet,
		c.apiURL(fmt.Sprintf("/external/integrations/thirdparty/getaccesspointstatusbyclient/%d", clientID)),
		nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RegisterUserTraffic records ISP user traffic data (download/upload counters).
//
// POST /external/integrations/thirdparty/isp/authentication_contracts/traffic
func (c *Client) RegisterUserTraffic(ctx context.Context, req RegisterUserTrafficRequest) (*TrafficResponse, error) {
	var result TrafficResponse
	if err := c.doJSON(ctx, http.MethodPost,
		c.apiURL("/external/integrations/thirdparty/isp/authentication_contracts/traffic"),
		req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
