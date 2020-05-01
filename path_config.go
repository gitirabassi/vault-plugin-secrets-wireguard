package wireguard

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	configPath = "config"
)

func (b *backend) pathConfig() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: configPath,
			Fields: map[string]*framework.FieldSchema{
				"port": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Port of the wireguard server.",
					Default:     51820,
				},
				"public_endpoint": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Address of the wireguard server.",
					Required:    true,
				},
				"client_persistent_keepalive": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "",
					Default:     25,
				},
				"save_config": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "",
					Default:     true,
				},
				"post_up_script": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: ".",
					Default:     "sysctl -w net.ipv4.ip_forward=1; iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE",
				},
				"post_down_script": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: ".",
					Default:     "iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE",
				},
				"server_cidr": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: ".",
					Default:     "10.29.0.1/24",
				},
				"webhook_address": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: ".",
				},
				"webhook_secret": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: ".",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.configCreateUpdateOperation,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.configCreateUpdateOperation,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.configReadOperation,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.configDeleteOperation,
				},
			},
			HelpSynopsis:    configHelpSynopsis,
			HelpDescription: configHelpDescription,
		},
	}
}

func (b *backend) configCreateUpdateOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	err := fieldData.Validate()
	if err != nil {
		return nil, errors.New("Failing validation: " + err.Error())
	}
	// Host must be provided all the time
	host := fieldData.Get("public_endpoint").(string)
	if host == "" {
		return nil, errors.New("public_endpoint is required")
	}
	// _, err = url.Parse(host)
	// if err != nil {
	// 	return nil, errors.New("host is not formatted correctly: " + err.Error())
	// }

	port := fieldData.Get("port").(int)
	config := &config{
		Port:           port,
		PublicEndpoint: host,
	}

	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Respond with a 204.
	return nil, nil
}

func (b *backend) configReadOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(ctx, configPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	config := &config{}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}
	var configMap map[string]interface{}
	err = mapstructure.Decode(config, &configMap)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: configMap,
	}
	return resp, nil
}

func (b *backend) configDeleteOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configPath); err != nil {
		return nil, err
	}
	return nil, nil
}

// config represent a wireguard/config
type config struct {
	Port                      int    `mapstructure:"port" `
	PrivateKey                string `mapstructure:"-"`
	PublicKey                 string `mapstructure:"publickey"`
	PublicEndpoint            string `mapstructure:"public_endpoint"`
	ClientPersistentKeepalive string `mapstructure:"client_persistent_keepalive"`
}

const configHelpSynopsis = `
Configure the Wireguard secret engine plugin.
`

const configHelpDescription = `
`
