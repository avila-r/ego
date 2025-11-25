package maps_test

import (
	"testing"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/maps"
	"github.com/stretchr/testify/assert"
)

func Test_Put(t *testing.T) {
	type Case struct {
		name     string
		entries  map[string]int
		expected map[string]int
	}

	cases := []Case{
		{"put single entry", map[string]int{"a": 1}, map[string]int{"a": 1}},
		{"put multiple entries", map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyHashMap[string, int]()
			for k, v := range c.entries {
				m.Put(k, v)
			}
			assert.Equal(t, c.expected, m.Elements())
		})
	}
}

func Test_Get(t *testing.T) {
	type Case struct {
		name     string
		init     map[string]int
		key      string
		expected int
		success  bool
	}

	cases := []Case{
		{"get existing key", map[string]int{"a": 10}, "a", 10, true},
		{"get non-existing key", map[string]int{"a": 10}, "b", 0, false},
		{"get from empty map", map[string]int{}, "x", 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyHashMap[string, int]()
			for k, v := range c.init {
				m.Put(k, v)
			}
			val, ok := m.Get(c.key)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_PutIfAbsent(t *testing.T) {
	type Case struct {
		name     string
		init     map[string]int
		key      string
		value    int
		expected map[string]int
		success  bool
	}

	cases := []Case{
		{"put if key absent", map[string]int{}, "a", 1, map[string]int{"a": 1}, true},
		{"do not put if key exists", map[string]int{"a": 1}, "a", 2, map[string]int{"a": 1}, false},
		{"put new key when others exist", map[string]int{"a": 1}, "b", 2, map[string]int{"a": 1, "b": 2}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyHashMap[string, int]()
			for k, v := range c.init {
				m.Put(k, v)
			}
			ok := m.PutIfAbsent(c.key, c.value)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, m.Elements())
		})
	}
}

func Test_Delete(t *testing.T) {
	type Case struct {
		name     string
		init     map[string]int
		key      string
		expected map[string]int
	}

	cases := []Case{
		{"delete existing key", map[string]int{"a": 1, "b": 2}, "a", map[string]int{"b": 2}},
		{"delete non-existing key", map[string]int{"a": 1}, "b", map[string]int{"a": 1}},
		{"delete from empty map", map[string]int{}, "x", map[string]int{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyHashMap[string, int]()
			for k, v := range c.init {
				m.Put(k, v)
			}
			m.Delete(c.key)
			assert.Equal(t, c.expected, m.Elements())
		})
	}
}

func Test_ContainsKey(t *testing.T) {
	type Case struct {
		name     string
		init     map[string]int
		key      string
		expected bool
	}

	cases := []Case{
		{"contains existing key", map[string]int{"a": 1}, "a", true},
		{"does not contain key", map[string]int{"a": 1}, "b", false},
		{"empty map", map[string]int{}, "x", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyHashMap[string, int]()
			for k, v := range c.init {
				m.Put(k, v)
			}
			assert.Equal(t, c.expected, m.ContainsKey(c.key))
		})
	}
}

func Test_ContainsValue(t *testing.T) {
	type Case struct {
		name     string
		init     map[string]int
		value    int
		expected bool
	}

	cases := []Case{
		{"contains existing value", map[string]int{"a": 10, "b": 20}, 10, true},
		{"does not contain value", map[string]int{"a": 10}, 20, false},
		{"empty map", map[string]int{}, 5, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyHashMap[string, int]()
			for k, v := range c.init {
				m.Put(k, v)
			}
			assert.Equal(t, c.expected, m.ContainsValue(c.value))
		})
	}
}

func Test_Len(t *testing.T) {
	type Case struct {
		name     string
		init     map[string]int
		expected int
	}

	cases := []Case{
		{"empty map", map[string]int{}, 0},
		{"three entries", map[string]int{"a": 1, "b": 2, "c": 3}, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := maps.EmptyHashMap[string, int]()
			for k, v := range c.init {
				m.Put(k, v)
			}
			assert.Equal(t, c.expected, m.Len())
		})
	}
}

func Test_IsEmpty(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	assert.True(t, m.IsEmpty())
	m.Put("a", 1)
	assert.False(t, m.IsEmpty())
}

func Test_Clear(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Clear()
	assert.Equal(t, 0, m.Len())
	assert.False(t, m.ContainsKey("a"))
	assert.False(t, m.ContainsKey("b"))
}

func Test_Filter(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
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
}

func Test_Clone(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)

	cloned := m.Clone()

	assert.Equal(t, m.Len(), cloned.Len())
	assert.Equal(t, m.Elements(), cloned.Elements())

	// Ensure independence
	cloned.Put("c", 3)
	assert.False(t, m.ContainsKey("c"))
	assert.True(t, cloned.ContainsKey("c"))
}

func Test_KeySlice(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)

	keys := m.KeySlice()
	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, "a")
	assert.Contains(t, keys, "b")
}

func Test_ValueSlice(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	m.Put("a", 10)
	m.Put("b", 20)

	values := m.ValueSlice()
	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, 10)
	assert.Contains(t, values, 20)
}

func Test_ToSlice(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)

	entries := m.ToSlice()
	assert.Equal(t, 2, len(entries))

	// Check that both entries exist (order doesn't matter)
	hasA := false
	hasB := false
	for _, entry := range entries {
		if entry.Key == "a" && entry.Value == 1 {
			hasA = true
		}
		if entry.Key == "b" && entry.Value == 2 {
			hasB = true
		}
	}
	assert.True(t, hasA)
	assert.True(t, hasB)
}

func Test_Keys(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)

	keys := m.Keys()
	assert.Equal(t, 2, keys.Size())
}

func Test_Values(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	m.Put("a", 10)
	m.Put("b", 20)

	values := m.Values()
	assert.Equal(t, 2, values.Size())
}

func Test_Entries(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
	m.Put("a", 1)

	entries := m.Entries()
	assert.Equal(t, 1, entries.Size())

	var entry collection.Entry[string, int]
	entries.ForEach(func(e collection.Entry[string, int]) {
		entry = e
	})
	assert.Equal(t, "a", entry.Key)
	assert.Equal(t, 1, entry.Value)
}

func Test_Iterator(t *testing.T) {
	m := maps.EmptyHashMap[string, int]()
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

	// Verify we got all 3 entries (order may vary for HashMap)
	assert.Equal(t, 3, len(entries))

	// Verify all expected entries are present
	keys := make(map[string]int)
	for _, entry := range entries {
		keys[entry.Key] = entry.Value
	}

	assert.Equal(t, 1, keys["a"])
	assert.Equal(t, 2, keys["b"])
	assert.Equal(t, 3, keys["c"])
}
