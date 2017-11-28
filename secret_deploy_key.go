package main

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	SecretDeployKeyType = "deploy_key"
)

func secretDeployKey(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretDeployKeyType,
		Fields: map[string]*framework.FieldSchema{
			"deploy_key": &framework.FieldSchema{
				Type: framework.TypeString,
			},
		},
		Revoke: b.secretDeployKeyRevoke,
	}
}

func (b *backend) secretDeployKeyRevoke(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyID, ok := req.Secret.InternalData["deploy_key_id"].(string)
	if !ok {
		return nil, nil
	}

	projectID, ok := req.Secret.InternalData["project_id"].(string)
	if !ok {
		return nil, nil
	}

	conf, err := readConfig(req.Storage)
	if err != nil {
		return nil, err
	}

	g, err := gitlabClient(conf.Address, conf.Token)

	if err != nil {
		return nil, err
	}

	if _, err := g.delDeployKey(projectID, keyID, b); err != nil {
		return nil, err
	}

	return nil, nil
}
