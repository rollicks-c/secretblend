package secretblend

import (
	"encoding/json"
	"fmt"
	"strings"
)

type injector[T any] struct {
}

type visitor func(key string, value interface{}) (interface{}, error)

func (i injector[T]) injectSecrets(subject T) (T, error) {

	// transform to generic map
	flat, err := i.toFlat(subject)
	if err != nil {
		return subject, err
	}

	// visit all nodes
	if err := i.visitNode(flat, i.processNode); err != nil {
		return subject, err
	}

	// transform back to subject
	injectedSubject, err := i.fromFlat(flat)
	if err != nil {
		return subject, err
	}

	return *injectedSubject, nil

}

func (i injector[T]) toFlat(subject T) (map[string]interface{}, error) {

	raw, err := json.Marshal(subject)
	if err != nil {
		return nil, err
	}

	var flat map[string]interface{}
	if err = json.Unmarshal(raw, &flat); err != nil {
		return nil, err
	}

	return flat, nil
}

func (i injector[T]) fromFlat(flat map[string]interface{}) (*T, error) {
	raw, err := json.Marshal(flat)
	if err != nil {
		return nil, err
	}

	var subject T
	if err = json.Unmarshal(raw, &subject); err != nil {
		return nil, err
	}

	return &subject, nil
}

func (i injector[T]) visitNode(item map[string]interface{}, visitor visitor) error {

	for key, value := range item {
		node, isNode := value.(map[string]interface{})
		if isNode {
			if err := i.visitNode(node, visitor); err != nil {
				return err
			}
		} else {
			newValue, err := visitor(key, value)
			if err != nil {
				return err
			}
			item[key] = newValue
		}
	}
	return nil
}

func (i injector[T]) processNode(key string, valueRaw interface{}) (interface{}, error) {

	// skip none-strings
	value, ok := valueRaw.(string)
	if !ok {
		return valueRaw, nil
	}

	// apply global injectors
	for _, gp := range globalProviders {
		processedValue, err := gp.LoadSecret(value)
		if err != nil {
			return nil, err
		}
		value = processedValue
	}

	// extract protocol
	parts := strings.Split(value, "://")
	if len(parts) != 2 {
		return value, nil
	}
	proto := protocol(fmt.Sprintf("%s://", parts[0]))
	secretURI := parts[1]

	// find provider
	provider, ok := protocolProviders[proto]
	if !ok {
		return value, nil
	}

	// inject secret
	secret, err := provider.LoadSecret(secretURI)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
