package wireguard

import "github.com/hashicorp/vault/sdk/framework"

func (b *backend) secretClientConfiguration() *framework.Secret {
	return &framework.Secret{
		Type: "authorization_token",
		Fields: map[string]*framework.FieldSchema{
			"conf": {
				Type:        framework.TypeString,
				Description: "Wireguard client configuration",
			},
		},
		// Renew:  b.authTokenRenew,
		// Revoke: b.authTokenRevoke,
	}
}

// func (b *backend) secretKeyRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
// 	resp, err := b.verifySecretServiceKeyExists(ctx, req)
// 	if err != nil {
// 		return resp, err
// 	}
// 	if resp == nil {
// 		resp = &logical.Response{}
// 	}
// 	cfg, err := getConfig(ctx, req.Storage)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if cfg == nil {
// 		cfg = &config{}
// 	}

// 	resp.Secret = req.Secret
// 	resp.Secret.TTL = cfg.TTL
// 	resp.Secret.MaxTTL = cfg.MaxTTL
// 	return resp, nil
// }

// func (b *backend) verifySecretServiceKeyExists(ctx context.Context, req *logical.Request) (*logical.Response, error) {
// 	keyName, ok := req.Secret.InternalData["key_name"]
// 	if !ok {
// 		return nil, fmt.Errorf("invalid secret, internal data is missing key name")
// 	}

// 	rsName, ok := req.Secret.InternalData["role_set"]
// 	if !ok {
// 		return nil, fmt.Errorf("invalid secret, internal data is missing role set name")
// 	}

// 	bindingSum, ok := req.Secret.InternalData["role_set_bindings"]
// 	if !ok {
// 		return nil, fmt.Errorf("invalid secret, internal data is missing role set checksum")
// 	}

// 	// Verify role set was not deleted.
// 	rs, err := getRoleSet(rsName.(string), ctx, req.Storage)
// 	if err != nil {
// 		return logical.ErrorResponse(fmt.Sprintf("could not find role set '%v' for secret", rsName)), nil
// 	}

// 	// Verify role set bindings have not changed since secret was generated.
// 	if rs.bindingHash() != bindingSum.(string) {
// 		return logical.ErrorResponse(fmt.Sprintf("role set '%v' bindings were updated since secret was generated, cannot renew", rsName)), nil
// 	}

// 	// Verify service account key still exists.
// 	iamAdmin, err := b.IAMAdminClient(req.Storage)
// 	if err != nil {
// 		return logical.ErrorResponse("could not confirm key still exists in GCP"), nil
// 	}
// 	if k, err := iamAdmin.Projects.ServiceAccounts.Keys.Get(keyName.(string)).Do(); err != nil || k == nil {
// 		return logical.ErrorResponse(fmt.Sprintf("could not confirm key still exists in GCP: %v", err)), nil
// 	}
// 	return nil, nil
// }

// func (b *backend) secretKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
// 	keyNameRaw, ok := req.Secret.InternalData["key_name"]
// 	if !ok {
// 		return nil, fmt.Errorf("secret is missing key_name internal data")
// 	}

// 	iamAdmin, err := b.IAMAdminClient(req.Storage)
// 	if err != nil {
// 		return logical.ErrorResponse(err.Error()), nil
// 	}

// 	_, err = iamAdmin.Projects.ServiceAccounts.Keys.Delete(keyNameRaw.(string)).Do()
// 	if err != nil && !isGoogleAccountKeyNotFoundErr(err) {
// 		return logical.ErrorResponse(fmt.Sprintf("unable to delete service account key: %v", err)), nil
// 	}

// 	return nil, nil
// }
