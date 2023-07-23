package powerbiapi

import (
	"fmt"
	"net/url"
)

/*
Add Datasource User
Grants or updates the permissions required to use the specified data source for the specified user.

Create Datasource
Creates a new data source on the specified on-premises gateway.

Delete Datasource
Deletes the specified data source from the specified gateway.

Delete Datasource User
Removes the specified user from the specified data source.

Get Datasource
Returns the specified data source from the specified gateway.

Get Datasources
Returns a list of data sources from the specified gateway.

Get Datasource Status
Checks the connectivity status of the specified data source from the specified gateway.

Get Datasource Users
Returns a list of users who have access to the specified data source.

Get Gateway
Returns the specified gateway.

Get Gateways
Returns a list of gateways for which the user is an admin.

Update Datasource
Updates the credentials of the specified data source from the specified gateway.
*/

type AddDatasouceUserRequest struct {
	DatasourceAccessRight string `json:"datasourceAccessRight,omitempty"`
	DisplayName           string `json:"displayName,omitempty"`
	EmailAddress          string `json:"emailAddress,omitempty"`
	GroupUserAccessRight  string `json:"groupUserAccessRight,omitempty"`
	Identifier            string `json:"identifier,omitempty"`
	PrincipalType         string `json:"principalType,omitempty"`
}

type CreateDatasourceRequest struct {
	ConnectionDetails string `json:"connectionDetails,omitempty"`
	//Value UpdateRefreshScheduleInGroupRequestValue `json:"value"`
	Value          CredentialDetails `json:"value"`
	DataSourceName string            `json:"dataSourceName,omitempty"`
	DataSourceType string            `json:"dataSourceType,omitempty"`
}

type CredentialDetails struct {
	CredentialType              string `json:"credentialType,omitempty"`
	Credentials                 string `json:"credentials,omitempty"`
	EncryptedConnection         string `json:"encryptedConnection,omitempty"`
	EncryptionAlgorithm         string `json:"encryptionAlgorithm,omitempty"`
	PrivacyLevel                string `json:"privacyLevel,omitempty"`
	UseCallerAADIdentity        *bool  `json:"useCallerAADIdentity,omitempty"`
	UseEndUserOAuth2Credentials *bool  `json:"useEndUserOAuth2Credentials,omitempty"`
}

// GetGatewayasResponse represents the response from the GetGatways API, an array of GatwayItems
type GetGatewaysResponse struct {
	Value []GetGatewaysResponseItem
}

type GetGatewaysResponseItem struct {
	ID                string
	Name              string
	Type              string
	PublicKey         GetGatewaysResponseItemGatewayPublicKey
	GatewayStatus     string
	GatewayAnnotation string
}

type GetGatewaysResponseItemGatewayPublicKey struct {
	exponent string
	modulus  string
}

type GetDatasourcesResponse struct {
	Value []GetDatasourcesResponseItem
}

type GetDatasourcesResponseItem struct {
	connectionDetails string
	//credentialDetails  needs investigation  > useEndUserOAuth2Credentials	boolean
	credentialType string
	datasourceName string
	datasourceType string
	gatewayId      string
	ID             string
}

type GetDatasourceStatusResponse struct {
	status string //todo: is a Json string
}

type GetDatasourceUsersResponse struct {
	Value []GetDatasourceUsersResponseItem
}

type GetDatasourceUsersResponseItem struct {
	datasourceAccessRight string
	displayName           string
	emailAddress          string
	identifier            string
	principalType         string
	profile               string
}

// CreateDatasource creates new datasource
func (client *Client) CreateDatasource(gatewayId string, request CreateDatasourceRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources", url.PathEscape(gatewayId))
	err := client.doJSON("POST", url, &request, nil)
	return err
}

// Grants or updates the permissions required to use the specified data source.
func (client *Client) DeleteDatasource(gatewayId string, datasourceId string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s", url.PathEscape(gatewayId), url.PathEscape(datasourceId))
	err := client.doJSON("DELETE", url, nil, nil)

	return err
}

// Grants or updates the permissions required to use the specified data source for the specified user.
func (client *Client) DeleteDatasourceUser(gatewayId string, datasourceId string, emailAdress string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s/users/%s", url.PathEscape(gatewayId), url.PathEscape(datasourceId), url.PathEscape(emailAdress))
	err := client.doJSON("DELETE", url, nil, nil)

	return err
}

// Removes the specified user from the specified data source.
func (client *Client) AddDatasourceUser(gatewayId string, datasourceId string, request AddDatasouceUserRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s/users", url.PathEscape(gatewayId), url.PathEscape(datasourceId))
	err := client.doJSON("POST", url, &request, nil)

	return err
}

// Returns a list of gateways for which the user is an admin.
func (client *Client) GetGateways() (*GetGatewaysResponse, error) {

	var respObj GetGatewaysResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways")
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// Returns the specified gateway.
func (client *Client) GetGateway(gatewayId string) (*GetGatewaysResponseItem, error) {

	var respObj GetGatewaysResponseItem
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s", url.PathEscape(gatewayId))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// Returns a list of data sources from the specified gateway.
func (client *Client) GetDatasources(gatewayId string) (*GetDatasourcesResponse, error) {

	var respObj GetDatasourcesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources", url.PathEscape(gatewayId))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// Returns the specified data source from the specified gateway.
func (client *Client) GetDatasource(gatewayId string, datasourceId string) (*GetDatasourcesResponseItem, error) {

	var respObj GetDatasourcesResponseItem
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s", url.PathEscape(gatewayId), url.PathEscape(datasourceId))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// Checks the connectivity status of the specified data source from the specified gateway.
func (client *Client) GetDatasourceStatus(gatewayId string, datasourceId string) (*GetDatasourceStatusResponse, error) {

	var respObj GetDatasourceStatusResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s/status", url.PathEscape(gatewayId), url.PathEscape(datasourceId))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

// Returns a list of users who have access to the specified data source.
func (client *Client) GetDatasourceUsers(gatewayId string, datasourceId string) (*GetDatasourceUsersResponse, error) {

	var respObj GetDatasourceUsersResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s/users", url.PathEscape(gatewayId), url.PathEscape(datasourceId))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}
