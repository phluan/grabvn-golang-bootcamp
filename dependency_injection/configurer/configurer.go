package configurer

import (
	"errors"
	httpclient "github.com/phluan/GrabGoTrainingWeek5Assignment/http_client"
	serializer "github.com/phluan/GrabGoTrainingWeek5Assignment/serializer"
)

type Configurer interface {
	HTTPClient() httpclient.HTTPClient
	Serializer() serializer.Serializer
}

type ConfigurerImpl struct {
	httpClient httpclient.HTTPClient
	serializer serializer.Serializer
}

func (impl *ConfigurerImpl) HTTPClient() httpclient.HTTPClient {
	return impl.httpClient
}

func (impl *ConfigurerImpl) Serializer() serializer.Serializer {
	return impl.serializer
}

type Option func(configure *ConfigurerImpl)

func WithHttpClient(httpClient httpclient.HTTPClient) Option {
	return func(configurer *ConfigurerImpl) {
		configurer.httpClient = httpClient
	}
}

func WithSerializer(serializer serializer.Serializer) Option {
	return func(configurer *ConfigurerImpl) {
		configurer.serializer = serializer
	}
}

func New(options ...Option) (Configurer, error) {
	configurer := &ConfigurerImpl{}

	for _, o := range options {
		o := o
		o(configurer)
	}

	if configurer.httpClient == nil {
		return nil, errors.New("missing http client")
	}

	if configurer.serializer == nil {
		return nil, errors.New("missing serializer")
	}

	return configurer, nil
}
