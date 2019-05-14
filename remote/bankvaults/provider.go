package bankvaults

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"

	"github.com/banzaicloud/bank-vaults/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/sagikazarmark/viperx/remote"
)

// nolint: gochecknoinits
func init() {
	remote.RegisterConfigProvider("bankvaults", NewConfigProvider())
}

// ConfigProvider implements reads configuration from Hashicorp Vault using Banzai Cloud Bank Vaults client.
type ConfigProvider struct{}

// NewConfigProvider returns a new ConfigProvider.
func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{}
}

func (p ConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	endpoint := rp.Endpoint()
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse provider endpoint")
	}

	query := u.Query()
	u.RawQuery = ""

	config := api.DefaultConfig()
	config.Address = u.String()
	rawClient, err := api.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create raw vault api client")
	}

	rawClient.SetToken(query.Get("token"))

	client, err := vault.NewClientFromRawClient(
		rawClient,
		vault.ClientRole(query.Get("role")),
		vault.ClientAuthPath(query.Get("authPath")),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create vault api client")
	}
	defer client.Close() // We close the client here to stop the unnecessary token renewal

	secret, err := client.RawClient().Logical().Read(rp.Path())
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
	return nil, errors.New("watch is not implemented for the bankvaults config provider")
}

func (p ConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	panic("watch channel is not implemented for the bankvaults config provider")
}
