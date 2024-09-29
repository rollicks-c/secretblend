package secretblend

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testSubject struct {
	Key1 string
	Key2 int
	Key3 bool
	Key4 string
}
type mockProvider struct {
	values map[string]string
}

func (m mockProvider) LoadSecret(uri string) (string, error) {
	if val, ok := m.values[uri]; ok {
		return val, nil
	}
	return "injected", nil
}

func TestNoChange(t *testing.T) {

	sub := testSubject{
		Key1: "value1",
		Key2: 2,
		Key3: true,
		Key4: "value4",
	}
	AddProvider(mockProvider{}, "mock://")
	subInjected, err := Inject(&sub)
	assert.NoError(t, err)
	assert.Equal(t, "value1", subInjected.Key1)
	assert.Equal(t, 2, subInjected.Key2)
	assert.Equal(t, true, subInjected.Key3)
	assert.Equal(t, "value4", subInjected.Key4)

}

func TestInject(t *testing.T) {

	sub := testSubject{
		Key1: "mock://value1",
		Key2: 2,
		Key3: true,
		Key4: "mock value4",
	}
	AddProvider(mockProvider{}, "mock://")
	subInjected, err := Inject(&sub)
	assert.NoError(t, err)
	assert.Equal(t, "injected", subInjected.Key1)
	assert.Equal(t, 2, subInjected.Key2)
	assert.Equal(t, true, subInjected.Key3)
	assert.Equal(t, "mock value4", subInjected.Key4)

}

func TestLookup(t *testing.T) {

	sub := testSubject{
		Key1: "mock://value1",
		Key2: 2,
		Key3: true,
		Key4: "mock://value4",
	}
	AddProvider(mockProvider{
		values: map[string]string{
			"value1": "secret1",
			"value4": "secret4",
		},
	}, "mock://")
	subInjected, err := Inject(&sub)
	assert.NoError(t, err)
	assert.Equal(t, "secret1", subInjected.Key1)
	assert.Equal(t, 2, subInjected.Key2)
	assert.Equal(t, true, subInjected.Key3)
	assert.Equal(t, "secret4", subInjected.Key4)

}
