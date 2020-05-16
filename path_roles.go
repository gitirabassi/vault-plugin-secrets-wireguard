package wireguard

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	rolesStoragePath = "roles/"
)

type roleEntry struct {
	Name        string   `json:"name" mapstructure:"name"`
	Permissions []string `json:"permissions" mapstructure:"permissions"`
	OrgID       string   `json:"org_id" mapstructure:"org_id"`
	TTL         int      `json:"ttl" mapstructure:"ttl"`
	MaxTTL      int      `json:"max_ttl" mapstructure:"max_ttl"`
}

func (b *backend) pathRoles() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: rolesStoragePath + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the role.",
				},
				"permissions": {
					Type:        framework.TypeStringSlice,
					Description: "List of permissions to give the token",
				},
				"org_id": {
					Type:        framework.TypeString,
					Description: "Organization ID in which to create the token",
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Default lease for generated credentials. If not set or set to 0, will use system default.",
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Maximum time a service principal. If not set or set to 0, will use system default.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathRoleCreateUpdate,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathRoleCreateUpdate,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathRoleRead,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathRoleDelete,
				},
			},
			HelpSynopsis:    roleHelpSyn,
			HelpDescription: roleHelpDesc,
		},
		{
			Pattern: "roles/?",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathRoleList,
				},
			},
			HelpSynopsis:    roleListHelpSyn,
			HelpDescription: roleListHelpDesc,
		},
	}
}

func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	b.Lock()
	defer b.Unlock()

	err := req.Storage.Delete(ctx, rolesStoragePath+name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	entry, err := req.Storage.Get(ctx, rolesStoragePath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	var roleMap map[string]interface{}
	err = mapstructure.Decode(result, &roleMap)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: roleMap,
	}, nil
}

func (b *backend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := data.Validate()
	if err != nil {
		return nil, errors.New("Failing validation: " + err.Error())
	}
	name := data.Get("name").(string)

	b.Lock()
	defer b.Unlock()

	orgID := data.Get("org_id").(string)
	if orgID == "" {
		return logical.ErrorResponse("OrgID is a required field to manage a influxdb v2 role"), nil
	}

	permissions := data.Get("permissions").([]string)
	if len(permissions) == 0 {
		return logical.ErrorResponse("permissions is a required field to manage a influxdb v2 role"), nil
	}

	ttl := data.Get("ttl").(int)
	maxTTL := data.Get("max_ttl").(int)

	roleEntry := &roleEntry{
		Name:        name,
		Permissions: permissions,
		OrgID:       orgID,
		TTL:         ttl,
		MaxTTL:      maxTTL,
	}
	entry, err := logical.StorageEntryJSON(rolesStoragePath+name, roleEntry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Respond with a 204.
	return nil, nil
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	rolesets, err := req.Storage.List(ctx, rolesStoragePath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(rolesets), nil
}

const roleHelpSyn = "Manage the Vault roles used to Manage access to Influxdb2."
const roleHelpDesc = `
Roles allow to define which permissions a given token will have

`

const roleListHelpSyn = "asdas"
const roleListHelpDesc = `
asdasda
asdasd
`
