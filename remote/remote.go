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

// SetErrorHandler sets the error handler of the global config provider registry.
func SetErrorHandler(errorHandler ErrorHandler) {
	r.SetErrorHandler(errorHandler)
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

// ErrorHandler handles an error occurred in a remote config provider.
type ErrorHandler interface {
	Handle(err error)
}

// ConfigProviderRegistry acts as a registry for remote container providers.
type ConfigProviderRegistry struct {
	configProviders map[string]ConfigProvider
	errorHandler    ErrorHandler
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

// SetErrorHandler sets the error handler of the registry.
func (r *ConfigProviderRegistry) SetErrorHandler(errorHandler ErrorHandler) {
	r.errorHandler = errorHandler
}

func (r *ConfigProviderRegistry) handleError(err error) {
	if r.errorHandler != nil {
		r.errorHandler.Handle(err)
	}
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
		r.handleError(err)

		return nil, err
	}

	re, err := configProvider.Get(rp)
	if err != nil {
		r.handleError(err)

		return nil, err
	}

	return re, nil
}

// Watch implements the ConfigProvider interface.
func (r *ConfigProviderRegistry) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	configProvider, err := r.getConfigProvider(rp)
	if err != nil {
		r.handleError(err)

		return nil, err
	}

	re, err := configProvider.Watch(rp)
	if err != nil {
		r.handleError(err)

		return nil, err
	}

	return re, nil
}

// WatchChannel implements the ConfigProvider interface.
func (r *ConfigProviderRegistry) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	configProvider, err := r.getConfigProvider(rp)
	if err != nil {
		panic(err)
	}

	return configProvider.WatchChannel(rp)
}
