package wireguard

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	serverCredsPath = "server-creds/"
)

func (b *backend) pathServerCreds() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: serverCredsPath + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the server.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathServerCredsRead,
				},
			},
			HelpSynopsis:    pathServerCredsReadHelpSyn,
			HelpDescription: pathServerCredsReadHelpDesc,
		},
	}
}

func (b *backend) pathServerCredsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	server, err := req.Storage.Get(ctx, serverStoragePath+name)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return logical.ErrorResponse("unknown server: %s", name), nil
	}

	result := &serverEntry{}
	if err := server.DecodeJSON(result); err != nil {
		return nil, err
	}
	if result == nil {
		return logical.ErrorResponse("server missconfigures (try deleting/readding it): %s", name), nil
	}
	// Getting config to embed "host" in the response
	configEntry, err := req.Storage.Get(ctx, configPath)
	if err != nil {
		return logical.ErrorResponse("influxdb config is missing (get): %s", err), nil
	}
	if configEntry == nil {
		return logical.ErrorResponse("influxdb config is missing (nil response): %s", err), nil
	}
	config := &config{}
	if err := configEntry.DecodeJSON(config); err != nil {
		return logical.ErrorResponse("influxdb config is missing (decode json): %s", err), nil
	}
	if config == nil {
		return logical.ErrorResponse("influxdb config is missing (decoded json is nil): %s", err), nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"host":   "asdasdas",
			"org_id": "asdasds",
			"token":  "sasd",
		},
	}
	if result.TTL != 0 {
		resp.Secret.TTL = time.Duration(result.TTL) * time.Second
	}
	if result.MaxTTL != 0 {
		resp.Secret.MaxTTL = time.Duration(result.MaxTTL) * time.Second
	}
	return resp, nil
}

const pathServerCredsReadHelpSyn = `
Request Influxdb credentials for a certain server. These credentials are
rotated periodically.`

const pathServerCredsReadHelpDesc = `
This path reads influxdb v2 credentials for a certain server. 
`
