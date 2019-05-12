package vault

import (
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func TestConfigProvider(t *testing.T) {
	address := os.Getenv("VAULT_ADDR")
	if address == "" {
		t.Skip("VAULT_ADDR environment variable not found")
	}

	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		t.Skip("VAULT_TOKEN environment variable not found")
	}

	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(token)

	remoteConfig := viper.RemoteConfig
	defer func() {
		viper.RemoteConfig = remoteConfig
	}()

	supportedRemoteProviders := viper.SupportedRemoteProviders
	defer func() {
		viper.SupportedRemoteProviders = supportedRemoteProviders
	}()

	viper.RemoteConfig = NewConfigProvider()
	viper.SupportedRemoteProviders = append(viper.SupportedRemoteProviders, "vault")

	_, err = client.Logical().Write("secret/data/hello", map[string]interface{}{
		"data": map[string]interface{}{
			"database": map[string]interface{}{
				"user": "root",
				"pass": "root",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse(address)
	if err != nil {
		t.Fatal(err)
	}

	u.RawQuery = "token=" + token

	v := viper.New()
	_ = v.AddRemoteProvider("vault", u.String(), "secret/data/hello")
	v.SetConfigType("json")

	err = v.ReadRemoteConfig()
	if err != nil {
		t.Fatal(err)
	}

	if v.GetString("database.user") != "root" || v.GetString("database.pass") != "root" {
		t.Errorf("failed to read secrets from vault")
	}

	_, err = client.Logical().Delete("secret/metadata/hello")
	if err != nil {
		t.Fatal(err)
	}
}
