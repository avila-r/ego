package maps_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/maps"
)

func TestEqual(t *testing.T) {
	type Map map[string]int

	type Case struct {
		name string
		l, r Map
		exp  bool
	}

	cases := []Case{
		{"empty", Map{}, Map{}, true},
		{"equal", Map{"a": 1}, Map{"a": 1}, true},
		{"different value", Map{"a": 1}, Map{"a": 2}, false},
		{"different key", Map{"a": 1}, Map{"b": 1}, false},
		{"different size", Map{"a": 1}, Map{"a": 1, "b": 2}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.exp, maps.Equal(c.l, c.r))
		})
	}
}

func TestEqualBy(t *testing.T) {
	type Map map[string]int

	type Case struct {
		name string
		l, r Map
		eq   func(int, int) bool
		exp  bool
	}

	always := func(a, b int) bool { return true }
	absEq := func(a, b int) bool { return a == b || a == -b }

	cases := []Case{
		{"abs equal", Map{"a": 1}, Map{"a": -1}, absEq, true},
		{"not abs equal", Map{"a": 1}, Map{"a": 3}, absEq, false},
		{"always", Map{"a": 100}, Map{"a": -999}, always, true},
		{"key mismatch", Map{"a": 1}, Map{"b": 1}, always, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.exp, maps.EqualBy(c.l, c.r, c.eq))
		})
	}
}

func TestClone(t *testing.T) {
	type Map map[string]int

	src := Map{"a": 1, "b": 2}
	clone := maps.Clone(src)

	assert.Equal(t, src, clone)
	assert.NotSame(t, &src, &clone)
}

func TestCopy(t *testing.T) {
	type Map map[string]int

	dst := Map{"x": 9}
	src := Map{"a": 1, "b": 2}

	maps.Copy(dst, src)

	assert.Equal(t, Map{"x": 9, "a": 1, "b": 2}, dst)
}

func TestDeleteIf(t *testing.T) {
	type Map map[string]int

	type Case struct {
		name string
		init Map
		del  func(string, int) bool
		exp  Map
	}

	cases := []Case{
		{
			"delete even values",
			Map{"a": 1, "b": 2, "c": 4},
			func(_ string, v int) bool { return v%2 == 0 },
			Map{"a": 1},
		},
		{
			"delete nothing",
			Map{"a": 1},
			func(_ string, _ int) bool { return false },
			Map{"a": 1},
		},
		{
			"delete all",
			Map{"a": 1, "b": 2},
			func(_ string, _ int) bool { return true },
			Map{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			maps.DeleteIf(c.init, c.del)
			assert.Equal(t, c.exp, c.init)
		})
	}
}

func TestEntries(t *testing.T) {
	type Map map[string]int

	m := Map{"a": 1, "b": 2}
	entries := maps.Entries(m).Elements()

	assert.ElementsMatch(t,
		[]collection.Entry[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
		},
		entries,
	)
}

func TestKeys(t *testing.T) {
	type Map map[string]int

	m := Map{"a": 1, "b": 2}
	keys := maps.Keys(m).Elements()

	assert.ElementsMatch(t, []string{"a", "b"}, keys)
}

func TestValues(t *testing.T) {
	type Map map[string]int

	m := Map{"a": 1, "b": 2}

	values := maps.Values(m).Elements()

	assert.ElementsMatch(t, []int{1, 2}, values)
}

func TestIter(t *testing.T) {
	type Map map[string]int

	m := Map{
		"a": 1,
		"b": 2,
	}

	var items []collection.Entry[string, int]

	it := maps.Iter(m)
	for it.HasNext() {
		items = append(items, it.Next())
	}

	assert.ElementsMatch(t,
		[]collection.Entry[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
		},
		items,
	)
}

func TestFrom(t *testing.T) {
	m := maps.From(maps.Map[string, int]{
		"a": 1,
		"b": 2,
	})

	assert.ElementsMatch(t,
		[]collection.Entry[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
		},
		m.ToSlice(),
	)
}

func TestOf(t *testing.T) {
	m := maps.Of(
		collection.Entry[string, int]{Key: "a", Value: 1},
		collection.Entry[string, int]{Key: "b", Value: 2},
	)

	assert.ElementsMatch(t,
		[]collection.Entry[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
		},
		m.ToSlice(),
	)
}

func TestNewEmpty(t *testing.T) {
	m := maps.New[string, int]()
	assert.Equal(t, 0, m.Len())

	m2 := maps.Empty[string, int]()
	assert.Equal(t, 0, m2.Len())
}

func TestLinkedFrom(t *testing.T) {
	m := maps.LinkedFrom(map[string]int{
		"a": 1,
		"b": 2,
	})

	assert.Equal(t,
		[]collection.Entry[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
		},
		m.ToSlice(),
	)
}

func TestLinkedOf(t *testing.T) {
	m := maps.LinkedOf(
		collection.Entry[string, int]{Key: "a", Value: 1},
		collection.Entry[string, int]{Key: "b", Value: 2},
	)

	assert.Equal(t,
		[]collection.Entry[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
		},
		m.ToSlice(),
	)
}
