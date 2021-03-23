package vault

import (
	"bytes"
	"encoding/json"
	"io"

	"emperror.dev/errors"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"

	"github.com/sagikazarmark/viperx/remote"
)

// nolint: gochecknoinits
func init() {
	remote.RegisterConfigProvider("vault", NewConfigProvider())
}

// ConfigProvider implements reads configuration from Hashicorp Vault.
type ConfigProvider struct {
	clients map[string]*api.Client
}

// NewConfigProvider returns a new ConfigProvider.
func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{
		clients: make(map[string]*api.Client),
	}
}

func (p ConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	// TODO: Initialize the client Once
	client, ok := p.clients[rp.Endpoint()]
	if !ok {
		endpoint := rp.Endpoint()

		c, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create vault api client")
		}

		client = c
		p.clients[endpoint] = c
	}

	secret, err := client.Logical().Read(rp.Path())
	if err != nil {
		return nil, errors.WrapIf(err, "failed to read secret")
	}

	if secret == nil {
		return nil, errors.Errorf("source not found: %s", rp.Path())
	}

	if secret.Data == nil && secret.Warnings != nil {
		return nil, errors.Errorf("source: %s errors: %v", rp.Path(), secret.Warnings)
	}

	secretData := secret.Data
	if secretV2, found := secret.Data["data"].(map[string]interface{}); found {
		secretData = secretV2
	}
	b, err := json.Marshal(secretData)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to json encode secret")
	}

	return bytes.NewReader(b), nil
}

func (p ConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	b, err := p.Get(rp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (p ConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	panic("watch channel is not implemented for the vault config provider")
}
