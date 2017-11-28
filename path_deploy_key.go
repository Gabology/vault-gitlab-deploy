package main

import (
	"strconv"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathDeployKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: framework.GenericNameRegex("project_id") + "/deploykeys",
		Fields: map[string]*framework.FieldSchema{
			"project_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Id of the project",
			},
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The public SSH key to create a deploy key for",
			},
			"can_push": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Can deploy key push to the project's repository",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathDeployKeyWrite,
		},
	}
}

func (b *backend) pathDeployKeyWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	conf, err := readConfig(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	g, err := gitlabClient(conf.Address, conf.Token)

	if err != nil {
		return nil, err
	}

	projectID := d.Get("project_id").(string)
	key := d.Get("key").(string)
	push := d.Get("can_push").(bool)

	deployKey, err := g.deployKey(projectID, key, push)

	if err != nil {
		return nil, err
	}

	s := b.Secret(SecretDeployKeyType).Response(map[string]interface{}{
		"deploy_key": deployKey.Key,
	}, map[string]interface{}{
		"deploy_key_id": strconv.Itoa(deployKey.ID),
		"project_id":    projectID,
	})
	s.Secret.Renewable = false
	s.Secret.TTL = time.Duration(conf.TTL) * time.Second

	return s, nil
}
