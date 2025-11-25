package maps_test

import (
	"testing"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/maps"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	m := maps.New[string, int]()
	assert.NotNil(t, m)
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Len())
}

func Test_Empty(t *testing.T) {
	m := maps.Empty[string, int]()
	assert.NotNil(t, m)
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Len())
}

func Test_Of(t *testing.T) {
	type Case struct {
		name     string
		entries  []collection.Entry[string, int]
		expected map[string]int
	}

	cases := []Case{
		{
			"create empty map with no entries",
			[]collection.Entry[string, int]{},
			map[string]int{},
		},
		{
			"create map with single entry",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
			map[string]int{"a": 1},
		},
		{
			"create map with multiple entries",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}},
			map[string]int{"a": 1, "b": 2, "c": 3},
		},
		{
			"create map with duplicate keys uses last value",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "a", Value: 2}},
			map[string]int{"a": 2},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.Of(c.entries...)
			assert.Equal(t, len(c.expected), m.Len())
			assert.Equal(t, c.expected, m.Elements())
		})
	}
}

func Test_NewLinked(t *testing.T) {
	m := maps.NewLinked[string, int]()
	assert.NotNil(t, m)
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Len())
}

func Test_EmptyLinked(t *testing.T) {
	m := maps.EmptyLinked[string, int]()
	assert.NotNil(t, m)
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Len())
}

func Test_OfLinked(t *testing.T) {
	type Case struct {
		name     string
		entries  []collection.Entry[string, int]
		expected []collection.Entry[string, int]
	}

	cases := []Case{
		{
			"create empty linked map with no entries",
			[]collection.Entry[string, int]{},
			[]collection.Entry[string, int]{},
		},
		{
			"create linked map with single entry",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
		},
		{
			"create linked map with multiple entries preserves order",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}},
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}},
		},
		{
			"create linked map with duplicate keys uses last value keeps first position",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "a", Value: 10}},
			[]collection.Entry[string, int]{{Key: "a", Value: 10}, {Key: "b", Value: 2}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.OfLinked(c.entries...)
			assert.Equal(t, len(c.expected), m.Len())
			assert.Equal(t, c.expected, m.ToSlice())
		})
	}
}

func Test_Of_CanUseMap(t *testing.T) {
	m := maps.Of(
		collection.Entry[string, int]{Key: "key1", Value: 100},
		collection.Entry[string, int]{Key: "key2", Value: 200},
	)

	val, ok := m.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	m.Put("key3", 300)
	assert.Equal(t, 3, m.Len())
}

func Test_OfLinked_CanUseMap(t *testing.T) {
	m := maps.OfLinked(
		collection.Entry[string, int]{Key: "key1", Value: 100},
		collection.Entry[string, int]{Key: "key2", Value: 200},
	)

	val, ok := m.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	m.Put("key3", 300)
	assert.Equal(t, 3, m.Len())

	// Verify order
	expected := []collection.Entry[string, int]{
		{Key: "key1", Value: 100},
		{Key: "key2", Value: 200},
		{Key: "key3", Value: 300},
	}
	assert.Equal(t, expected, m.ToSlice())
}
