package maps_test

import (
	"testing"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/maps"
	"github.com/stretchr/testify/assert"
)

func TestLinked_Put(t *testing.T) {
	type Case struct {
		name     string
		entries  []collection.Entry[string, int]
		expected []collection.Entry[string, int]
	}

	cases := []Case{
		{
			"put single entry",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
		},
		{
			"put multiple entries maintains order",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}},
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}},
		},
		{
			"overwrite existing key maintains position",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "a", Value: 10}},
			[]collection.Entry[string, int]{{Key: "a", Value: 10}, {Key: "b", Value: 2}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyLinkedHashMap[string, int]()
			for _, entry := range c.entries {
				m.Put(entry.Key, entry.Value)
			}
			assert.Equal(t, c.expected, m.ToSlice())
		})
	}
}

func TestLinked_Get(t *testing.T) {
	type Case struct {
		name     string
		init     []collection.Entry[string, int]
		key      string
		expected int
		success  bool
	}

	cases := []Case{
		{"get existing key", []collection.Entry[string, int]{{Key: "a", Value: 10}}, "a", 10, true},
		{"get non-existing key", []collection.Entry[string, int]{{Key: "a", Value: 10}}, "b", 0, false},
		{"get from empty map", []collection.Entry[string, int]{}, "x", 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyLinkedHashMap[string, int]()
			for _, entry := range c.init {
				m.Put(entry.Key, entry.Value)
			}
			val, ok := m.Get(c.key)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, val)
		})
	}
}

func TestLinked_PutIfAbsent(t *testing.T) {
	type Case struct {
		name     string
		init     []collection.Entry[string, int]
		key      string
		value    int
		expected []collection.Entry[string, int]
		success  bool
	}

	cases := []Case{
		{
			"put if key absent",
			[]collection.Entry[string, int]{},
			"a", 1,
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
			true,
		},
		{
			"do not put if key exists",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
			"a", 2,
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
			false,
		},
		{
			"put new key when others exist",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
			"b", 2,
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}},
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyLinkedHashMap[string, int]()
			for _, entry := range c.init {
				m.Put(entry.Key, entry.Value)
			}
			ok := m.PutIfAbsent(c.key, c.value)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, m.ToSlice())
		})
	}
}

func TestLinked_Delete(t *testing.T) {
	type Case struct {
		name     string
		init     []collection.Entry[string, int]
		key      string
		expected []collection.Entry[string, int]
	}

	cases := []Case{
		{
			"delete first element",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}},
			"a",
			[]collection.Entry[string, int]{{Key: "b", Value: 2}, {Key: "c", Value: 3}},
		},
		{
			"delete middle element",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}},
			"b",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "c", Value: 3}},
		},
		{
			"delete last element",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}},
			"c",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}},
		},
		{
			"delete non-existing key",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
			"b",
			[]collection.Entry[string, int]{{Key: "a", Value: 1}},
		},
		{
			"delete from empty map",
			[]collection.Entry[string, int]{},
			"x",
			[]collection.Entry[string, int]{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyLinkedHashMap[string, int]()
			for _, entry := range c.init {
				m.Put(entry.Key, entry.Value)
			}
			m.Delete(c.key)
			assert.Equal(t, c.expected, m.ToSlice())
		})
	}
}

func TestLinked_ContainsKey(t *testing.T) {
	type Case struct {
		name     string
		init     []collection.Entry[string, int]
		key      string
		expected bool
	}

	cases := []Case{
		{"contains existing key", []collection.Entry[string, int]{{Key: "a", Value: 1}}, "a", true},
		{"does not contain key", []collection.Entry[string, int]{{Key: "a", Value: 1}}, "b", false},
		{"empty map", []collection.Entry[string, int]{}, "x", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyLinkedHashMap[string, int]()
			for _, entry := range c.init {
				m.Put(entry.Key, entry.Value)
			}
			assert.Equal(t, c.expected, m.ContainsKey(c.key))
		})
	}
}

func TestLinked_ContainsValue(t *testing.T) {
	type Case struct {
		name     string
		init     []collection.Entry[string, int]
		value    int
		expected bool
	}

	cases := []Case{
		{"contains existing value", []collection.Entry[string, int]{{Key: "a", Value: 10}, {Key: "b", Value: 20}}, 10, true},
		{"does not contain value", []collection.Entry[string, int]{{Key: "a", Value: 10}}, 20, false},
		{"empty map", []collection.Entry[string, int]{}, 5, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyLinkedHashMap[string, int]()
			for _, entry := range c.init {
				m.Put(entry.Key, entry.Value)
			}
			assert.Equal(t, c.expected, m.ContainsValue(c.value))
		})
	}
}

func TestLinked_Len(t *testing.T) {
	type Case struct {
		name     string
		init     []collection.Entry[string, int]
		expected int
	}

	cases := []Case{
		{"empty map", []collection.Entry[string, int]{}, 0},
		{"three entries", []collection.Entry[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}}, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyLinkedHashMap[string, int]()
			for _, entry := range c.init {
				m.Put(entry.Key, entry.Value)
			}
			assert.Equal(t, c.expected, m.Len())
		})
	}
}

func TestLinked_IsEmpty(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	assert.True(t, m.IsEmpty())
	m.Put("a", 1)
	assert.False(t, m.IsEmpty())
}

func TestLinked_Clear(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Clear()
	assert.Equal(t, 0, m.Len())
	assert.False(t, m.ContainsKey("a"))
	assert.False(t, m.ContainsKey("b"))
	assert.Equal(t, []collection.Entry[string, int]{}, m.ToSlice())
}

func TestLinked_Filter(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 10)
	m.Put("b", 20)
	m.Put("c", 30)

	filtered := m.Filter(func(k string, v int) bool {
		return v > 15
	})

	assert.Equal(t, 2, filtered.Len())
	assert.True(t, filtered.ContainsKey("b"))
	assert.True(t, filtered.ContainsKey("c"))
	assert.False(t, filtered.ContainsKey("a"))

	// Verify order is maintained
	expected := []collection.Entry[string, int]{{Key: "b", Value: 20}, {Key: "c", Value: 30}}
	assert.Equal(t, expected, filtered.ToSlice())
}

func TestLinked_Clone(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	cloned := m.Clone()

	assert.Equal(t, m.Len(), cloned.Len())
	assert.Equal(t, m.ToSlice(), cloned.ToSlice())

	// Ensure independence
	cloned.Put("d", 4)
	assert.False(t, m.ContainsKey("d"))
	assert.True(t, cloned.ContainsKey("d"))
}

func TestLinked_KeySlice(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	keys := m.KeySlice()
	expected := []string{"a", "b", "c"}
	assert.Equal(t, expected, keys)
}

func TestLinked_ValueSlice(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 10)
	m.Put("b", 20)
	m.Put("c", 30)

	values := m.ValueSlice()
	expected := []int{10, 20, 30}
	assert.Equal(t, expected, values)
}

func TestLinked_ToSlice(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	entries := m.ToSlice()
	expected := []collection.Entry[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}
	assert.Equal(t, expected, entries)
}

func TestLinked_Keys(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)

	keys := m.Keys()
	assert.Equal(t, 2, keys.Size())
}

func TestLinked_Values(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 10)
	m.Put("b", 20)

	values := m.Values()
	assert.Equal(t, 2, values.Size())
}

func TestLinked_Entries(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)

	entries := m.Entries()
	assert.Equal(t, 2, entries.Size())
}

func TestLinked_InsertionOrder(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()

	// Insert in specific order
	m.Put("z", 26)
	m.Put("a", 1)
	m.Put("m", 13)
	m.Put("b", 2)

	// Verify order is maintained (not alphabetical)
	expected := []collection.Entry[string, int]{
		{Key: "z", Value: 26},
		{Key: "a", Value: 1},
		{Key: "m", Value: 13},
		{Key: "b", Value: 2},
	}
	assert.Equal(t, expected, m.ToSlice())
}

func TestLinked_UpdatePreservesOrder(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	// Update middle element
	m.Put("b", 20)

	expected := []collection.Entry[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 20},
		{Key: "c", Value: 3},
	}
	assert.Equal(t, expected, m.ToSlice())
}

func TestLinked_Elements(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	elements := m.Elements()

	expected := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	assert.Equal(t, expected, elements)
	assert.Equal(t, 3, len(elements))
}

func TestLinked_Iterator(t *testing.T) {
	m := maps.EmptyLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	iter := m.Iterator()

	assert.NotNil(t, iter)

	// Collect entries from iterator
	var entries []collection.Entry[string, int]
	for iter.HasNext() {
		entry := iter.Next()
		entries = append(entries, entry)
	}

	expected := []collection.Entry[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}

	assert.Equal(t, expected, entries)
}
