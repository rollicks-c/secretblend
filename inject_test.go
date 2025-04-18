package secretblend

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testSubjectSimple struct {
	Key1 string
	Key2 int
	Key3 bool
	Key4 string
}

type testSubjectNested struct {
	Key1 string
	Key2 int
	Key3 testSubjectSimple
	Key4 testSubjectSimple
}

type testSubjectLists struct {
	Key1 string
	Key2 int
	Key3 []testSubjectSimple
	Key4 []testSubjectSimple
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

	sub := testSubjectSimple{
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

	sub := testSubjectSimple{
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

func TestInjectNested(t *testing.T) {

	sub := testSubjectNested{
		Key1: "mock://value1",
		Key2: 2,
		Key3: testSubjectSimple{
			Key1: "mock://value3",
			Key2: 3,
			Key3: true,
			Key4: "mock value4",
		},
		Key4: testSubjectSimple{
			Key1: "mock://value5",
			Key2: 5,
			Key3: true,
			Key4: "mock value6",
		},
	}
	AddProvider(mockProvider{
		values: map[string]string{
			"value1": "secret1",
			"value3": "secret3",
			"value5": "secret5",
		},
	}, "mock://")
	subInjected, err := Inject(&sub)
	assert.NoError(t, err)
	assert.Equal(t, "secret1", subInjected.Key1)
	assert.Equal(t, 2, subInjected.Key2)
	assert.Equal(t, true, subInjected.Key3.Key3)
	assert.Equal(t, "mock value4", subInjected.Key3.Key4)
	assert.Equal(t, "secret1", subInjected.Key1)
	assert.Equal(t, 2, subInjected.Key2)
	assert.Equal(t, true, subInjected.Key3.Key3)
	assert.Equal(t, "mock value4", subInjected.Key3.Key4)
	assert.Equal(t, "secret3", subInjected.Key3.Key1)
	assert.Equal(t, 3, subInjected.Key3.Key2)
	assert.Equal(t, true, subInjected.Key3.Key3)
	assert.Equal(t, "mock value4", subInjected.Key3.Key4)
	assert.Equal(t, "secret5", subInjected.Key4.Key1)
	assert.Equal(t, 5, subInjected.Key4.Key2)
	assert.Equal(t, true, subInjected.Key4.Key3)
	assert.Equal(t, "mock value6", subInjected.Key4.Key4)
}

func TestLookup(t *testing.T) {

	sub := testSubjectSimple{
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

func TestInjectSlices(t *testing.T) {

	sub := testSubjectLists{
		Key1: "mock://value1",
		Key2: 2,
		Key3: []testSubjectSimple{
			{
				Key1: "mock://value31.1",
				Key2: 32,
				Key3: true,
				Key4: "mock value34",
			},
			{
				Key1: "mock://value31.2",
				Key2: 32,
				Key3: true,
				Key4: "mock value34",
			},
		},
		Key4: []testSubjectSimple{
			{
				Key1: "mock://value41.1",
				Key2: 42,
				Key3: true,
				Key4: "mock value44",
			},
		},
	}
	AddProvider(mockProvider{
		values: map[string]string{
			"value1":    "secret1",
			"value31.1": "secret31.1",
			"value31.2": "secret31.2",
			"value41.1": "secret41.1",
		},
	}, "mock://")

	subInjected, err := Inject(&sub)
	assert.NoError(t, err)
	assert.Equal(t, "secret1", subInjected.Key1)
	assert.Equal(t, 2, subInjected.Key2)
	assert.Equal(t, "secret31.1", subInjected.Key3[0].Key1)
	assert.Equal(t, 32, subInjected.Key3[0].Key2)
	assert.Equal(t, true, subInjected.Key3[0].Key3)
	assert.Equal(t, "mock value34", subInjected.Key3[0].Key4)
	assert.Equal(t, "secret31.2", subInjected.Key3[1].Key1)
	assert.Equal(t, 32, subInjected.Key3[1].Key2)
	assert.Equal(t, true, subInjected.Key3[1].Key3)
	assert.Equal(t, "mock value34", subInjected.Key3[1].Key4)
	assert.Equal(t, "secret41.1", subInjected.Key4[0].Key1)
	assert.Equal(t, 42, subInjected.Key4[0].Key2)
	assert.Equal(t, true, subInjected.Key4[0].Key3)
	assert.Equal(t, "mock value44", subInjected.Key4[0].Key4)

}
