package vault

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
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
	client, ok := p.clients[rp.Endpoint()]
	if !ok {
		endpoint := rp.Endpoint()
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse provider endpoint")
		}

		query := u.Query()
		u.RawQuery = ""

		config := api.DefaultConfig()
		config.Address = u.String()
		c, err := api.NewClient(config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create vault api client")
		}

		c.SetToken(query.Get("token"))

		client = c
		p.clients[endpoint] = c
	}

	secret, err := client.Logical().Read(rp.Path())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read secret")
	}

	if secret == nil {
		return nil, errors.Errorf("source not found: %s", rp.Path())
	}

	if secret.Data == nil && secret.Warnings != nil {
		return nil, errors.Errorf("source: %s errors: %v", rp.Path(), secret.Warnings)
	}

	b, err := json.Marshal(secret.Data["data"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to json encode secret")
	}

	return bytes.NewReader(b), nil
}

func (p ConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return nil, errors.New("watch is not implemented for the vault config provider")
}

func (p ConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	panic("watch channel is not implemented for the vault config provider")
}
