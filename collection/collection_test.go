package collection_test

import (
	"testing"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	type Case struct {
		name     string
		items    []int
		expected []int
	}

	cases := []Case{
		{"create with no items", []int{}, []int{}},
		{"create with single item", []int{1}, []int{1}},
		{"create with multiple items", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.items...)
			assert.Equal(t, c.expected, col.Elements())
		})
	}
}

func Test_Empty(t *testing.T) {
	col := collection.Empty[int]()
	assert.NotNil(t, col)
	assert.True(t, col.IsEmpty())
	assert.Equal(t, 0, col.Size())
	assert.Equal(t, []int{}, col.Elements())
}

func Test_Sized(t *testing.T) {
	type Case struct {
		name     string
		size     int
		expected int
	}

	cases := []Case{
		{"create with size 0", 0, 0},
		{"create with size 10", 10, 0},
		{"create with size 100", 100, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Sized[int](c.size)
			assert.NotNil(t, col)
			assert.Equal(t, c.expected, col.Size())
			assert.True(t, col.IsEmpty())
		})
	}
}

func Test_Add(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		toAdd    []int
		expected []int
	}

	cases := []Case{
		{"add to empty collection", []int{}, []int{1}, []int{1}},
		{"add single element", []int{1}, []int{2}, []int{1, 2}},
		{"add multiple elements", []int{1}, []int{2, 3, 4}, []int{1, 2, 3, 4}},
		{"add to non-empty collection", []int{10, 20}, []int{30, 40}, []int{10, 20, 30, 40}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.initial...)
			col.Add(c.toAdd...)
			assert.Equal(t, c.expected, col.Elements())
		})
	}
}

func Test_Get(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		index    int
		expected int
		success  bool
	}

	cases := []Case{
		{"get valid index", []int{10, 20, 30}, 1, 20, true},
		{"get first element", []int{10, 20, 30}, 0, 10, true},
		{"get last element", []int{10, 20, 30}, 2, 30, true},
		{"get negative index", []int{10, 20}, -1, 0, false},
		{"get out of range positive", []int{10, 20}, 5, 0, false},
		{"get from empty collection", []int{}, 0, 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)
			val, ok := col.Get(c.index)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Remove(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		index    int
		expected []int
		success  bool
	}

	cases := []Case{
		{"remove first element", []int{1, 2, 3}, 0, []int{2, 3}, true},
		{"remove middle element", []int{10, 20, 30, 40}, 2, []int{10, 20, 40}, true},
		{"remove last element", []int{5, 6, 7}, 2, []int{5, 6}, true},
		{"remove only element", []int{42}, 0, []int{}, true},
		{"remove negative index", []int{1, 2}, -1, []int{1, 2}, false},
		{"remove out of range positive", []int{1, 2}, 5, []int{1, 2}, false},
		{"remove from empty collection", []int{}, 0, []int{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)
			ok := col.Remove(c.index)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, col.Elements())
		})
	}
}

func Test_Size(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		expected int
	}

	cases := []Case{
		{"empty collection", []int{}, 0},
		{"single element", []int{1}, 1},
		{"multiple elements", []int{1, 2, 3, 4, 5}, 5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)
			assert.Equal(t, c.expected, col.Size())
		})
	}
}

func Test_IsEmpty(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		expected bool
	}

	cases := []Case{
		{"empty collection", []int{}, true},
		{"non-empty collection", []int{1}, false},
		{"multiple elements", []int{1, 2, 3}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)
			assert.Equal(t, c.expected, col.IsEmpty())
		})
	}
}

func Test_Clear(t *testing.T) {
	type Case struct {
		name string
		init []int
	}

	cases := []Case{
		{"clear empty collection", []int{}},
		{"clear single element", []int{1}},
		{"clear multiple elements", []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)
			col.Clear()
			assert.Equal(t, 0, col.Size())
			assert.True(t, col.IsEmpty())
			assert.Equal(t, []int{}, col.Elements())
		})
	}
}

func Test_Elements(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		expected []int
	}

	cases := []Case{
		{"empty collection", []int{}, []int{}},
		{"single element", []int{1}, []int{1}},
		{"multiple elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)
			elements := col.Elements()
			assert.Equal(t, c.expected, elements)
		})
	}
}

func Test_Elements_ReturnsClone(t *testing.T) {
	col := collection.Of(1, 2, 3)
	elements := col.Elements()

	// Modify returned slice
	elements[0] = 999

	// Original collection should be unchanged
	val, _ := col.Get(0)
	assert.Equal(t, 1, val)
}

func Test_Clone(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		expected []int
	}

	cases := []Case{
		{"clone empty collection", []int{}, []int{}},
		{"clone single element", []int{1}, []int{1}},
		{"clone multiple elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)
			cloned := col.Clone()
			assert.Equal(t, c.expected, cloned.Elements())
			assert.Equal(t, col.Size(), cloned.Size())
		})
	}
}

func Test_Clone_Independence(t *testing.T) {
	col := collection.Of(1, 2, 3)
	cloned := col.Clone()

	// Modify original
	col.Add(4)

	// Cloned should be unchanged
	assert.Equal(t, 4, col.Size())
	assert.Equal(t, 3, cloned.Size())
	assert.Equal(t, []int{1, 2, 3}, cloned.Elements())

	// Modify cloned
	cloned.Add(99)

	// Original should be unchanged
	assert.Equal(t, []int{1, 2, 3, 4}, col.Elements())
	assert.Equal(t, []int{1, 2, 3, 99}, cloned.Elements())
}

func Test_ForEach(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		expected []int
	}

	cases := []Case{
		{"iterate empty collection", []int{}, []int{}},
		{"iterate single element", []int{1}, []int{1}},
		{"iterate multiple elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)

			result := slice.Empty[int]()
			col.ForEach(func(item int) {
				result = append(result, item)
			})

			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_ForEach_Mutation(t *testing.T) {
	col := collection.Of(1, 2, 3)
	sum := 0
	col.ForEach(func(item int) {
		sum += item
	})
	assert.Equal(t, 6, sum)
}

func Test_Iterator(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		expected []int
	}

	cases := []Case{
		{"iterate empty collection", []int{}, []int{}},
		{"iterate single element", []int{1}, []int{1}},
		{"iterate multiple elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.init...)
			iter := col.Iterator()
			assert.NotNil(t, iter)

			result := slice.Empty[int]()
			for iter.HasNext() {
				result = append(result, iter.Next())
			}

			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Collection_ComplexWorkflow(t *testing.T) {
	// Create collection
	col := collection.Empty[string]()
	assert.True(t, col.IsEmpty())

	// Add elements
	col.Add("a", "b", "c")
	assert.Equal(t, 3, col.Size())
	assert.False(t, col.IsEmpty())

	// Get elements
	val, ok := col.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "b", val)

	// Remove element
	ok = col.Remove(1)
	assert.True(t, ok)
	assert.Equal(t, 2, col.Size())
	assert.Equal(t, []string{"a", "c"}, col.Elements())

	// Clone
	cloned := col.Clone()
	cloned.Add("d")
	assert.Equal(t, 2, col.Size())
	assert.Equal(t, 3, cloned.Size())

	// Clear
	col.Clear()
	assert.True(t, col.IsEmpty())
	assert.Equal(t, 0, col.Size())
}
