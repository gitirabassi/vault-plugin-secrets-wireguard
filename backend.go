package wireguard

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config",
				"role/*",
			},
		},
		Paths: framework.PathAppend(
			b.pathConfig(),
			b.pathRoles(),
			b.pathServers(),
			b.pathCreds(),
			b.pathServerCreds(),
		),
		Secrets: []*framework.Secret{
			b.secretClientConfiguration(),
			b.secretServerConfiguration(),
		},
		BackendType: logical.TypeLogical,
	}
	return &b
}

type backend struct {
	*framework.Backend
	sync.RWMutex
}

const backendHelp = `
The Wireguard backend supports managing a wireguard server

After mounting this secret backend, configure it using the "wireguard/config" path.
`
