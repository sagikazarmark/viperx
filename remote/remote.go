package remote

import (
	"io"

	"github.com/spf13/viper"
)

// nolint: gochecknoinits
func init() {
	viper.RemoteConfig = r
}

// nolint: gochecknoglobals
var r = NewConfigProviderRegistry()

// RegisterConfigProvider registers a config provider in the global config provider registry.
func RegisterConfigProvider(provider string, configProvider ConfigProvider) {
	r.RegisterConfigProvider(provider, configProvider)

	AddSupportedRemoteProvider(provider)
}

// AddSupportedRemoteProvider adds a remote provider to the list of supported providers.
func AddSupportedRemoteProvider(provider string) {
	if !isStringInSlice(provider, viper.SupportedRemoteProviders) {
		viper.SupportedRemoteProviders = append(viper.SupportedRemoteProviders, provider)
	}
}

func isStringInSlice(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}

	return false
}

// ConfigProvider is the interface defined by Viper for remote config providers.
type ConfigProvider interface {
	Get(rp viper.RemoteProvider) (io.Reader, error)
	Watch(rp viper.RemoteProvider) (io.Reader, error)
	WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool)
}

// ConfigProviderRegistry acts as a registry for remote container providers.
type ConfigProviderRegistry struct {
	configProviders map[string]ConfigProvider
}

// NewConfigProviderRegistry returns a new ConfigProviderRegistry.
func NewConfigProviderRegistry() *ConfigProviderRegistry {
	return &ConfigProviderRegistry{
		configProviders: make(map[string]ConfigProvider),
	}
}

// RegisterConfigProvider registers a config provider in the registry.
func (r *ConfigProviderRegistry) RegisterConfigProvider(provider string, configProvider ConfigProvider) {
	r.configProviders[provider] = configProvider
}

func (r *ConfigProviderRegistry) getConfigProvider(rp viper.RemoteProvider) (ConfigProvider, error) {
	provider := rp.Provider()
	configProvider, ok := r.configProviders[provider]
	if !ok {
		return nil, viper.UnsupportedRemoteProviderError(provider)
	}

	return configProvider, nil
}

// Get implements the ConfigProvider interface.
func (r *ConfigProviderRegistry) Get(rp viper.RemoteProvider) (io.Reader, error) {
	configProvider, err := r.getConfigProvider(rp)
	if err != nil {
		return nil, err
	}

	return configProvider.Get(rp)
}

// Watch implements the ConfigProvider interface.
func (r *ConfigProviderRegistry) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	configProvider, err := r.getConfigProvider(rp)
	if err != nil {
		return nil, err
	}

	return configProvider.Watch(rp)
}

// WatchChannel implements the ConfigProvider interface.
func (r *ConfigProviderRegistry) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	configProvider, err := r.getConfigProvider(rp)
	if err != nil {
		panic(err)
	}

	return configProvider.WatchChannel(rp)
}
