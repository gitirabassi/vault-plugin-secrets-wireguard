package wireguard

import (
	"context"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	defaultLeaseTTLVal = time.Hour * 12
	maxLeaseTTLVal     = time.Hour * 24
)

func getBackend(throwsErr bool) (logical.Backend, logical.Storage) {
	config := &logical.BackendConfig{
		Logger: logging.NewVaultLogger(log.Error),

		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
		StorageView: &logical.InmemStorage{},
	}
	b := Backend()
	b.Setup(context.Background(), config)

	return b, config.StorageView
}
