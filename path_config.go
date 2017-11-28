package main

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Gitlab server address",
			},
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Personal access token for API calls",
			},
			"deploy_key_ttl": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "TTL in seconds for deploy keys (default 5 minutes)",
				Default:     300,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},
	}
}

func (b *backend) pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config", config{
		Address: d.Get("address").(string),
		Token:   d.Get("token").(string),
		TTL:     d.Get("deploy_key_ttl").(int),
	})

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func readConfig(storage logical.Storage) (*config, error) {
	entry, err := storage.Get("config")
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("Gitlab backend has not been configured. Please configure it at the '/config' endpoint first")
	}

	conf := &config{}
	entry.DecodeJSON(conf)

	return conf, nil
}

func (b *backend) pathConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	conf, err := readConfig(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": conf.Address,
			"token":   conf.Token,
			"ttl":     conf.TTL,
		},
	}, nil
}

type config struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	TTL     int    `json:"deploy_key_ttl"`
}
