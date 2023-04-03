package vault

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func TestConfigProvider(t *testing.T) {
	var secretPath = "secret/data/hello"

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

	// extract root path from the secret path
	mountPath := strings.Split(secretPath, "/")[0] + "/"

	if !mountExists(client, mountPath) {
		createMountOrDie(t, client, mountPath)
		defer deleteMount(client, mountPath)
	}

	_, err = client.Logical().Write(secretPath, map[string]interface{}{
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
	_ = v.AddRemoteProvider("vault", u.String(), secretPath)
	v.SetConfigType("json")

	err = v.ReadRemoteConfig()
	if err != nil {
		t.Fatal(err)
	}

	if v.GetString("database.user") != "root" || v.GetString("database.pass") != "root" {
		t.Errorf("failed to read secrets from vault")
	}

	_, err = client.Logical().Delete(secretPath)
	if err != nil {
		t.Fatal(err)
	}
}

func mountExists(client *api.Client, mountPath string) bool {
	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return false
	}

	_, ok := mounts[mountPath]
	return ok
}

func createMountOrDie(t *testing.T, client *api.Client, mountPath string) {
	err := client.Sys().Mount(mountPath, &api.MountInput{
		Type:        "kv",
		Description: "viperx auto tests",
		Config:      api.MountConfigInput{},
	})
	if err != nil {
		t.Fatalf("secret mount path not found for %s and failed to create it as a fixture: %v", mountPath, err)
	}
}

func deleteMount(client *api.Client, mountPath string) {
	_ = client.Sys().Unmount(mountPath)
}
