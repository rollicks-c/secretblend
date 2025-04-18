package secretblend

type Provider interface {
	LoadSecret(uri string) (string, error)
}

func AddProvider(provider Provider, protocols ...string) {
	for _, proto := range protocols {
		providerList[protocol(proto)] = provider
	}
}
func HasProvider(protocolName string) bool {
	_, isRegistered := providerList[protocol(protocolName)]
	return isRegistered
}

func Inject[T any](subject T) (T, error) {
	inject := injector[T]{}
	injectedSubject, err := inject.injectSecrets(subject)
	if err != nil {
		return subject, err
	}
	return injectedSubject, nil
}
