package main

import (
	"log"
	"os"

	"github.com/hashicorp/vault/logical/framework"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args)

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: Factory,
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		log.Fatal(err)
	}
}

func Factory(c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)

	if err := b.Setup(c); err != nil {
		return nil, err
	}

	return b, nil
}

type backend struct {
	*framework.Backend
}

func Backend(c *logical.BackendConfig) *backend {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Secrets: []*framework.Secret{
			secretDeployKey(&b),
		},
		Paths: []*framework.Path{
			pathConfig(&b),
			pathDeployKey(&b),
		},
	}

	return &b
}
