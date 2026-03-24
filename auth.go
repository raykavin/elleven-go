package elleven

import (
	"context"
	"net/http"
	"net/url"
)

// Auth (Autenticação)

// Default OAuth2 platform constants used by the ERP Elleven legacy flows.
const (
	DefaultLegacyClientID     = "synauth"
	DefaultLegacyClientSecret = "df956154024a425eb80f1a2fc12fef0c"
	DefaultLegacyGrantType    = "password"
	DefaultLegacyScope        = "syngw synpaygw offline_access"

	DefaultClientCredentialsGrantType = "client_credentials"
	DefaultClientCredentialsScope     = "syngw"

	DefaultRefreshTokenGrantType = "refresh_token"
)

// strOr returns val if non-empty, otherwise returns def.
func strOr(val, def string) string {
	if val != "" {
		return val
	}
	return def
}

// TokenResponse is returned by all authentication endpoints.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// LegacyAuthRequest holds credentials for the legacy password-based OAuth2 flow.
//
// The platform defaults (GrantType, Scope, ClientID, ClientSecret) are pre-filled
// with the ERP Elleven standard values. Override them only if your environment requires
// different values.
type LegacyAuthRequest struct {
	// GrantType defaults to "password".
	GrantType string
	// Scope defaults to "syngw synpaygw offline_access".
	Scope string
	// ClientID defaults to "synauth".
	ClientID string
	// ClientSecret defaults to the ERP Elleven platform secret.
	ClientSecret string
	// Username is the integration user's login name.
	Username string
	// Password is the integration user's password.
	Password string
	// SynData is the integration token from Suite / Settings / Parameters / Integration/Map.
	SynData string
}

// ClientCredentialsRequest holds credentials for the modern OAuth2 client_credentials flow.
//
// ClientID and ClientSecret are obtained from the integration user record in
// Settings / Users / "Client Id" and "Client Secret" fields.
// GrantType and Scope default to the ERP Elleven standard values.
type ClientCredentialsRequest struct {
	// GrantType defaults to "client_credentials".
	GrantType string
	// Scope defaults to "syngw".
	Scope string
	// ClientID is obtained from Settings / Users / Integration User / "Client Id".
	ClientID string
	// ClientSecret is obtained from Settings / Users / Integration User / "Client Secret".
	ClientSecret string
	// SynData is the integration token from Suite / Settings / Parameters / Integration/Map.
	SynData string
}

// RefreshTokenRequest holds the data needed to refresh an existing access token.
//
// GrantType, ClientID and ClientSecret default to the ERP Elleven legacy platform
// values. Override them only if your environment requires different values.
type RefreshTokenRequest struct {
	// GrantType defaults to "refresh_token".
	GrantType string
	// ClientID defaults to "synauth".
	ClientID string
	// ClientSecret defaults to the ERP Elleven platform secret.
	ClientSecret string
	// RefreshToken is the token obtained from a previous authentication response.
	RefreshToken string
}

// AuthenticateLegacy authenticates using the legacy password grant OAuth2 flow.
//
// The returned TokenResponse can be passed directly to client.SetTokenResponse.
func (c *Client) AuthenticateLegacy(ctx context.Context, req LegacyAuthRequest) (*TokenResponse, error) {
	form := url.Values{
		"grant_type":    {strOr(req.GrantType, DefaultLegacyGrantType)},
		"scope":         {strOr(req.Scope, DefaultLegacyScope)},
		"client_id":     {strOr(req.ClientID, DefaultLegacyClientID)},
		"client_secret": {strOr(req.ClientSecret, DefaultLegacyClientSecret)},
		"username":      {req.Username},
		"password":      {req.Password},
		"syndata":       {req.SynData},
	}

	var result TokenResponse
	if err := c.doForm(ctx, http.MethodPost, c.authURL("/connect/token"), form, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AuthenticateClientCredentials authenticates using the modern client_credentials
// OAuth2 flow.
//
// The returned TokenResponse can be passed directly to client.SetTokenResponse.
func (c *Client) AuthenticateClientCredentials(ctx context.Context, req ClientCredentialsRequest) (*TokenResponse, error) {
	form := url.Values{
		"grant_type":    {strOr(req.GrantType, DefaultClientCredentialsGrantType)},
		"scope":         {strOr(req.Scope, DefaultClientCredentialsScope)},
		"client_id":     {req.ClientID},
		"client_secret": {req.ClientSecret},
		"syndata":       {req.SynData},
	}

	var result TokenResponse
	if err := c.doForm(ctx, http.MethodPost, c.authURL("/connect/token"), form, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RefreshToken exchanges a refresh token for a new access token.
//
// The returned TokenResponse can be passed directly to client.SetTokenResponse.
func (c *Client) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*TokenResponse, error) {
	form := url.Values{
		"grant_type":    {strOr(req.GrantType, DefaultRefreshTokenGrantType)},
		"client_id":     {strOr(req.ClientID, DefaultLegacyClientID)},
		"client_secret": {strOr(req.ClientSecret, DefaultLegacyClientSecret)},
		"refresh_token": {req.RefreshToken},
	}

	var result TokenResponse
	if err := c.doForm(ctx, http.MethodPost, c.authURL("/connect/token"), form, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
