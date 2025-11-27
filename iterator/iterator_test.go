package iterator_test

import (
	"testing"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/iterator"
	"github.com/avila-r/ego/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		expected []int
	}

	cases := []Case{
		{"create with no elements", []int{}, []int{}},
		{"create with single element", []int{1}, []int{1}},
		{"create with multiple elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)
			assert.NotNil(t, it)

			result := slice.Empty[int]()
			for it.HasNext() {
				result = append(result, it.Next())
			}

			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_From(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		expected []int
	}

	cases := []Case{
		{"create from empty collection", []int{}, []int{}},
		{"create from single element", []int{1}, []int{1}},
		{"create from multiple elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := collection.Of(c.elements...)
			it := iterator.From(col)
			assert.NotNil(t, it)

			result := slice.Empty[int]()
			for it.HasNext() {
				result = append(result, it.Next())
			}

			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_HasNext(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		checks   []bool
	}

	cases := []Case{
		{"empty iterator", []int{}, []bool{false}},
		{"single element", []int{1}, []bool{true, false}},
		{"multiple elements", []int{1, 2, 3}, []bool{true, true, true, false}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)

			for _, expected := range c.checks {
				has := it.HasNext()
				assert.Equal(t, expected, has)
				if has {
					it.Next()
				}
			}
		})
	}
}

func Test_Next(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		expected []int
	}

	cases := []Case{
		{"iterate single element", []int{42}, []int{42}},
		{"iterate multiple elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"iterate in order", []int{10, 20, 30}, []int{10, 20, 30}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)

			result := slice.Empty[int]()
			for it.HasNext() {
				result = append(result, it.Next())
			}
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Next_Panics(t *testing.T) {
	it := iterator.Of[int]()
	assert.False(t, it.HasNext())
	assert.Panics(t, func() {
		it.Next()
	})
}

func Test_Next_PanicsAfterExhaustion(t *testing.T) {
	it := iterator.Of(1, 2, 3)

	it.Next()
	it.Next()
	it.Next()

	assert.False(t, it.HasNext())
	assert.Panics(t, func() {
		it.Next()
	})
}

func Test_Peek(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		peeks    []int
	}

	cases := []Case{
		{"peek first element", []int{1, 2, 3}, []int{1, 1, 1}},
		{"peek single element", []int{42}, []int{42, 42}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)

			for _, expected := range c.peeks {
				assert.Equal(t, expected, it.Peek())
			}
		})
	}
}

func Test_Peek_DoesNotAdvance(t *testing.T) {
	it := iterator.Of(1, 2, 3)

	// Peek multiple times
	assert.Equal(t, 1, it.Peek())
	assert.Equal(t, 1, it.Peek())
	assert.Equal(t, 1, it.Peek())

	// Next should still return first element
	assert.Equal(t, 1, it.Next())
	assert.Equal(t, 2, it.Peek())
}

func Test_Peek_Panics(t *testing.T) {
	it := iterator.Of[int]()
	assert.False(t, it.HasNext())
	assert.Panics(t, func() {
		it.Peek()
	})
}

func Test_Reset(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		consume  int
	}

	cases := []Case{
		{"reset after one element", []int{1, 2, 3}, 1},
		{"reset after all elements", []int{1, 2, 3}, 3},
		{"reset without consuming", []int{1, 2, 3}, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)

			// Consume some elements
			for i := 0; i < c.consume; i++ {
				it.Next()
			}

			// Reset
			it.Reset()

			// Should be able to iterate again from start
			result := slice.Empty[int]()
			for it.HasNext() {
				result = append(result, it.Next())
			}
			assert.Equal(t, c.elements, result)
		})
	}
}

func Test_Remaining(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		consume  int
		expected int
	}

	cases := []Case{
		{"no elements consumed", []int{1, 2, 3}, 0, 3},
		{"one element consumed", []int{1, 2, 3}, 1, 2},
		{"all elements consumed", []int{1, 2, 3}, 3, 0},
		{"empty iterator", []int{}, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)

			for i := 0; i < c.consume; i++ {
				it.Next()
			}

			assert.Equal(t, c.expected, it.Remaining())
		})
	}
}

func Test_Collect(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		consume  int
		expected []int
	}

	cases := []Case{
		{"collect all elements", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"collect after one consumed", []int{1, 2, 3}, 1, []int{2, 3}},
		{"collect after all consumed", []int{1, 2, 3}, 3, []int{}},
		{"collect empty iterator", []int{}, 0, []int{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)

			for i := 0; i < c.consume; i++ {
				it.Next()
			}

			result := it.Collect()
			assert.Equal(t, c.expected, result)

			// After collect, iterator should be exhausted
			assert.False(t, it.HasNext())
		})
	}
}

func Test_ForEach(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		expected []int
	}

	cases := []Case{
		{"iterate empty", []int{}, []int{}},
		{"iterate single element", []int{1}, []int{1}},
		{"iterate multiple elements", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)

			result := slice.Empty[int]()
			it.ForEach(func(elem int) {
				result = append(result, elem)
			})

			assert.Equal(t, c.expected, result)

			// After ForEach, iterator should be exhausted
			assert.False(t, it.HasNext())
		})
	}
}

func Test_ForEach_WithMutation(t *testing.T) {
	it := iterator.Of(1, 2, 3, 4, 5)

	sum := 0
	it.ForEach(func(elem int) {
		sum += elem
	})

	assert.Equal(t, 15, sum)
}

func Test_ForEach_PartialIteration(t *testing.T) {
	it := iterator.Of(1, 2, 3, 4, 5)

	// Consume first element
	it.Next()

	result := slice.Empty[int]()
	it.ForEach(func(elem int) {
		result = append(result, elem)
	})

	// Should only get remaining elements
	assert.Equal(t, []int{2, 3, 4, 5}, result)
}

func Test_Filter(t *testing.T) {
	type Case struct {
		name      string
		elements  []int
		predicate func(int) bool
		expected  []int
	}

	cases := []Case{
		{
			"filter even numbers",
			[]int{1, 2, 3, 4, 5, 6},
			func(n int) bool { return n%2 == 0 },
			[]int{2, 4, 6},
		},
		{
			"filter greater than 3",
			[]int{1, 2, 3, 4, 5},
			func(n int) bool { return n > 3 },
			[]int{4, 5},
		},
		{
			"filter none match",
			[]int{1, 2, 3},
			func(n int) bool { return n > 10 },
			[]int{},
		},
		{
			"filter all match",
			[]int{1, 2, 3},
			func(n int) bool { return n > 0 },
			[]int{1, 2, 3},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)
			filtered := it.Filter(c.predicate)

			result := slice.Empty[int]()
			for filtered.HasNext() {
				result = append(result, filtered.Next())
			}

			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Filter_ExhaustsOriginal(t *testing.T) {
	it := iterator.Of(1, 2, 3, 4, 5)

	filtered := it.Filter(func(n int) bool { return n%2 == 0 })

	// Original iterator should be exhausted
	assert.False(t, it.HasNext())

	// Filtered iterator should work
	assert.True(t, filtered.HasNext())
}

func Test_Map(t *testing.T) {
	type Case struct {
		name     string
		elements []int
		mapper   func(int) int
		expected []int
	}

	cases := []Case{
		{
			"map multiply by 2",
			[]int{1, 2, 3, 4, 5},
			func(n int) int { return n * 2 },
			[]int{2, 4, 6, 8, 10},
		},
		{
			"map add 10",
			[]int{1, 2, 3},
			func(n int) int { return n + 10 },
			[]int{11, 12, 13},
		},
		{
			"map to negative",
			[]int{1, 2, 3},
			func(n int) int { return -n },
			[]int{-1, -2, -3},
		},
		{
			"map empty iterator",
			[]int{},
			func(n int) int { return n * 2 },
			[]int{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			it := iterator.Of(c.elements...)
			mapped := iterator.Map(it, c.mapper)

			result := slice.Empty[int]()
			for mapped.HasNext() {
				result = append(result, mapped.Next())
			}

			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Map_TypeTransformation(t *testing.T) {
	it := iterator.Of(1, 2, 3)

	// Map int to string
	mapped := iterator.Map(it, func(n int) string {
		return string(rune('a' + n - 1))
	})

	var result []string
	for mapped.HasNext() {
		result = append(result, mapped.Next())
	}

	expected := []string{"a", "b", "c"}
	assert.Equal(t, expected, result)
}

func Test_Map_ExhaustsOriginal(t *testing.T) {
	it := iterator.Of(1, 2, 3)

	mapped := iterator.Map(it, func(n int) int { return n * 2 })

	// Original iterator should be exhausted
	assert.False(t, it.HasNext())

	// Mapped iterator should work
	assert.True(t, mapped.HasNext())
}

func Test_Iterator_ComplexWorkflow(t *testing.T) {
	// Create iterator
	it := iterator.Of(1, 2, 3, 4, 5)

	// Peek first element
	assert.Equal(t, 1, it.Peek())

	// Consume first element
	assert.Equal(t, 1, it.Next())

	// Check remaining
	assert.Equal(t, 4, it.Remaining())

	// Consume second element
	assert.Equal(t, 2, it.Next())

	// Collect rest
	rest := it.Collect()
	assert.Equal(t, []int{3, 4, 5}, rest)

	// Should be exhausted
	assert.False(t, it.HasNext())

	// Reset and iterate again
	it.Reset()
	assert.True(t, it.HasNext())
	assert.Equal(t, 5, it.Remaining())
}

func Test_Iterator_ChainedOperations(t *testing.T) {
	it := iterator.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// Filter even numbers
	filtered := it.Filter(func(n int) bool { return n%2 == 0 })

	// Map to multiply by 10
	mapped := iterator.Map(filtered, func(n int) int { return n * 10 })

	result := slice.Empty[int]()
	for mapped.HasNext() {
		result = append(result, mapped.Next())
	}

	expected := []int{20, 40, 60, 80, 100}
	assert.Equal(t, expected, result)
}

func Test_Iterator_MultipleResets(t *testing.T) {
	it := iterator.Of(1, 2, 3)

	// First iteration
	assert.Equal(t, 1, it.Next())
	it.Reset()

	// Second iteration
	assert.Equal(t, 1, it.Next())
	assert.Equal(t, 2, it.Next())
	it.Reset()

	// Third iteration
	result := slice.Empty[int]()
	for it.HasNext() {
		result = append(result, it.Next())
	}

	assert.Equal(t, []int{1, 2, 3}, result)
}
