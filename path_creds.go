package wireguard

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	credsPath = "creds/"
)

func (b *backend) pathCreds() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: credsPath + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the role.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathCredsRead,
				},
			},
			HelpSynopsis:    pathCredsReadHelpSyn,
			HelpDescription: pathCredsReadHelpDesc,
		},
	}
}

func (b *backend) pathCredsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	role, err := req.Storage.Get(ctx, rolesStoragePath+name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("unknown role: %s", name), nil
	}

	result := &roleEntry{}
	if err := role.DecodeJSON(result); err != nil {
		return nil, err
	}
	if result == nil {
		return logical.ErrorResponse("role missconfigures (try deleting/readding it): %s", name), nil
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
			"host":   "asdasds",
			"org_id": "asdasdas",
			"token":  "asdas",
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

const pathCredsReadHelpSyn = `
Request Influxdb credentials for a certain role. These credentials are
rotated periodically.`

const pathCredsReadHelpDesc = `
This path reads influxdb v2 credentials for a certain role. 
`
