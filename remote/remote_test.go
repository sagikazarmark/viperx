package remote

import (
	"bytes"
	"io"
	"testing"

	"github.com/spf13/viper"
)

type inMemoryConfigProvider struct{}

func (p *inMemoryConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	return bytes.NewReader([]byte(`{"key": "value"}`)), nil
}

func (p *inMemoryConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	panic("implement me")
}

func (p *inMemoryConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	panic("implement me")
}

func TestAddSupportedRemoteProvider(t *testing.T) {
	AddSupportedRemoteProvider("vault")
	defer viper.Reset()

	var containsProvider bool
	for _, provider := range viper.SupportedRemoteProviders {
		if provider == "vault" {
			containsProvider = true

			break
		}
	}

	if !containsProvider {
		t.Error("the list of supported remote providers was expected to contain \"vault\"")
	}
}

func TestRegisterConfigProvider(t *testing.T) {
	remoteConfig := viper.RemoteConfig
	defer func() {
		viper.RemoteConfig = remoteConfig
	}()

	supportedRemoteProviders := viper.SupportedRemoteProviders
	defer func() {
		viper.SupportedRemoteProviders = supportedRemoteProviders
	}()

	RegisterConfigProvider("inmemory", &inMemoryConfigProvider{})

	v := viper.New()

	err := v.AddRemoteProvider("inmemory", "inmemory", "inmemory")
	if err != nil {
		t.Fatal(err)
	}

	v.SetConfigType("json")

	err = v.ReadRemoteConfig()
	if err != nil {
		t.Fatal(err)
	}

	if got, want := v.GetString("key"), "value"; got != want {
		t.Errorf("remote config is not read: expected key not found\nactual:  %q\nexpected: %q", got, want)
	}
}

func TestConfigProviderRegistry_RegisterConfigProvider(t *testing.T) {
	registry := NewConfigProviderRegistry()

	configProvider := &inMemoryConfigProvider{}
	registry.RegisterConfigProvider("inmemory", configProvider)

	if registry.configProviders["inmemory"] != configProvider {
		t.Error("failed to register config provider")
	}
}

type errorHandler struct {
	err error
}

func (h *errorHandler) Handle(err error) {
	h.err = err
}

func TestConfigProviderRegistry_ErrorHandler(t *testing.T) {
	remoteConfig := viper.RemoteConfig
	defer func() {
		viper.RemoteConfig = remoteConfig
	}()

	supportedRemoteProviders := viper.SupportedRemoteProviders
	defer func() {
		viper.SupportedRemoteProviders = supportedRemoteProviders
	}()

	errorHandler := &errorHandler{}
	registry := NewConfigProviderRegistry()
	registry.SetErrorHandler(errorHandler)

	viper.RemoteConfig = registry
	AddSupportedRemoteProvider("inmemory")

	v := viper.New()

	err := v.AddRemoteProvider("inmemory", "inmemory", "inmemory")
	if err != nil {
		t.Fatal(err)
	}

	v.SetConfigType("json")

	err = v.ReadRemoteConfig()
	if err == nil {
		t.Fatal("expected an error")
	}

	if errorHandler.err == nil {
		t.Fatal("expected an error")
	}
}
