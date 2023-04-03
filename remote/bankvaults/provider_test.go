package bankvaults

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func TestConfigProvider(t *testing.T) {
	var secretPath = "secret/data1/hello"

	// extract root path from the secret path
	rootPath := strings.Split(secretPath, "/")[0] + "/"

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
	viper.SupportedRemoteProviders = append(viper.SupportedRemoteProviders, "bankvaults")

	// check the root mount for the secrets exists
	// create a mount if it doesn't exist and cleanup after the tests
	if !validMount(client, rootPath) {
		err = client.Sys().Mount(rootPath, &api.MountInput{
			Type:        "kv",
			Description: "viperx auto tests",
			Config:      api.MountConfigInput{},
		})
		if err != nil {
			t.Fatalf("secret mount path not found for %s and failed to create it as a fixture: %v", rootPath, err)
		}
		defer func() {
			err = client.Sys().Unmount(rootPath)
			if err != nil {
				t.Fatalf("failed to unmount the secret path %s: %v", rootPath, err)
			}
			fmt.Printf("unmounted the secret root %s", rootPath)
		}()
	}

	// check the kv version
	kvVersion, err := determineKVversion(client, rootPath)
	if err != nil {
		t.Fatal(err)
	}
	if kvVersion != 1 {
		t.Fatalf("kv version %d not supported", kvVersion)
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
	_ = v.AddRemoteProvider("bankvaults", u.String(), secretPath)
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

// validateMount checks the root mount and path are available since we rely on
func validMount(client *api.Client, secretPath string) bool {
	rootPath := strings.Split(secretPath, "/")[0] + "/"

	// check the root mount for the secrets exists
	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return false
	}
	for mountName, _ := range mounts {
		// fmt.Printf("mount %s: %+v\n", k, v.Type)
		if mountName == rootPath {
			return true
		}
	}

	// make sure it's readable for the clent
	scrt, err := client.Logical().Read(rootPath)
	if err != nil || scrt == nil {
		return false
	}
	return false
}

//
func determineKVversion(client *api.Client, secretPath string) (int, error) {
	// check the root mount for the secrets exists
	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return 0, err
	}
	for mountName, mountOutput := range mounts {
		if mountName == secretPath {
			fmt.Printf("mount %s output: %+v\n", mountName, mountOutput)
			if mountOutput.Type != "kv" {
				return 0, fmt.Errorf("mount %s is not a kv mount", mountName)
			}
			if mountOutput.Options["version"] == "2" {
				return 2, nil
			} else {
				return 1, nil
			}
		}
	}
	return 0, fmt.Errorf("could not determine KV version for %s", secretPath)
}
