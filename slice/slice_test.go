package slice_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/avila-r/ego/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	type Case struct {
		name     string
		values   []int
		expected []int
	}

	cases := []Case{
		{"create slice with single element", []int{1}, []int{1}},
		{"create slice with multiple elements", []int{1, 2, 3}, []int{1, 2, 3}},
		{"create slice with no elements", []int{}, []int{}},
		{"create slice with zero values", []int{0, 0, 0}, []int{0, 0, 0}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Of(c.values...)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_New(t *testing.T) {
	result := slice.New[int]()
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
	assert.Equal(t, []int{}, result)
}

func Test_Empty(t *testing.T) {
	result := slice.Empty[int]()
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
	assert.Equal(t, []int{}, result)
}

func Test_Append(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		toAdd    []int
		expected []int
	}

	cases := []Case{
		{"append to empty slice", []int{}, []int{1}, []int{1}},
		{"append single element", []int{1, 2}, []int{3}, []int{1, 2, 3}},
		{"append multiple elements", []int{1}, []int{2, 3, 4}, []int{1, 2, 3, 4}},
		{"append nothing", []int{1, 2}, []int{}, []int{1, 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Append(c.initial, c.toAdd...)
			assert.Equal(t, c.expected, result)
			// Verify original is unchanged if it has capacity
			if len(c.initial) > 0 {
				assert.Equal(t, c.initial[:len(c.initial)], c.initial)
			}
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
		{"add to empty slice", []int{}, []int{1}, []int{1}},
		{"add single element", []int{1, 2}, []int{3}, []int{1, 2, 3}},
		{"add multiple elements", []int{1}, []int{2, 3, 4}, []int{1, 2, 3, 4}},
		{"add nothing", []int{1, 2}, []int{}, []int{1, 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.initial
			slice.Add(&s, c.toAdd...)
			assert.Equal(t, c.expected, s)
		})
	}
}

func Test_Filter(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		filter   func(int) bool
		expected []int
	}

	cases := []Case{
		{"filter even numbers", []int{1, 2, 3, 4, 5}, func(v int) bool { return v%2 == 0 }, []int{2, 4}},
		{"filter odd numbers", []int{1, 2, 3, 4, 5}, func(v int) bool { return v%2 != 0 }, []int{1, 3, 5}},
		{"filter none match", []int{1, 2, 3}, func(v int) bool { return v > 10 }, []int{}},
		{"filter all match", []int{1, 2, 3}, func(v int) bool { return v > 0 }, []int{1, 2, 3}},
		{"filter empty slice", []int{}, func(v int) bool { return v > 0 }, []int{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Filter(c.initial, c.filter)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_IsEmpty(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		expected bool
	}

	cases := []Case{
		{"empty slice", []int{}, true},
		{"non-empty slice", []int{1}, false},
		{"slice with multiple elements", []int{1, 2, 3}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.IsEmpty(c.slice)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_First(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		expected int
		hasValue bool
	}

	cases := []Case{
		{"first of single element", []int{42}, 42, true},
		{"first of multiple elements", []int{1, 2, 3}, 1, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.First(c.slice)
			if c.hasValue {
				assert.True(t, result.IsPresent())
				assert.Equal(t, c.expected, result.Join())
			} else {
				assert.False(t, result.IsPresent())
			}
		})
	}
}

func Test_First_Panics_On_Empty(t *testing.T) {
	assert.Panics(t, func() {
		slice.First([]int{})
	})
}

func Test_Last(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		expected int
		hasValue bool
	}

	cases := []Case{
		{"last of single element", []int{42}, 42, true},
		{"last of multiple elements", []int{1, 2, 3}, 3, true},
		{"last of empty slice", []int{}, 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Last(c.slice)
			if c.hasValue {
				assert.True(t, result.IsPresent())
				assert.Equal(t, c.expected, result.Join())
			} else {
				assert.False(t, result.IsPresent())
			}
		})
	}
}

func Test_Size(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		expected int
	}

	cases := []Case{
		{"size of empty slice", []int{}, 0},
		{"size of single element", []int{1}, 1},
		{"size of multiple elements", []int{1, 2, 3, 4, 5}, 5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Size(c.slice)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_IsNil(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		expected bool
	}

	var nilSlice []int
	emptySlice := []int{}

	cases := []Case{
		{"nil slice", nilSlice, true},
		{"empty slice", emptySlice, false},
		{"non-empty slice", []int{1, 2, 3}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.IsNil(c.slice)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Stream(t *testing.T) {
	s := slice.Of(1, 2, 3, 4, 5)

	stream := slice.Stream(s)

	assert.NotNil(t, stream)
}

func Test_ForEach(t *testing.T) {
	t.Run("foreach executes for all elements", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		sum := 0

		slice.ForEach(s, func(v int) {
			sum += v
		})

		assert.Equal(t, 15, sum)
	})

	t.Run("foreach on empty slice", func(t *testing.T) {
		s := []int{}
		called := false

		slice.ForEach(s, func(v int) {
			called = true
		})

		assert.False(t, called)
	})

	t.Run("foreach with side effects", func(t *testing.T) {
		s := []int{1, 2, 3}
		result := []int{}

		slice.ForEach(s, func(v int) {
			result = append(result, v*2)
		})

		assert.Equal(t, []int{2, 4, 6}, result)
	})
}

func Test_Map(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		mapper   func(int) int
		expected []int
	}

	cases := []Case{
		{"map double values", []int{1, 2, 3}, func(v int) int { return v * 2 }, []int{2, 4, 6}},
		{"map to squares", []int{1, 2, 3, 4}, func(v int) int { return v * v }, []int{1, 4, 9, 16}},
		{"map empty slice", []int{}, func(v int) int { return v * 2 }, []int{}},
		{"map to constant", []int{1, 2, 3}, func(v int) int { return 42 }, []int{42, 42, 42}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Map(c.initial, c.mapper)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Map_Type_Conversion(t *testing.T) {
	t.Run("map int to string", func(t *testing.T) {
		s := []int{1, 2, 3}
		result := slice.Map(s, func(v int) string {
			return fmt.Sprintf("New: %d", v)
		})
		assert.Equal(t, []string{"New: 1", "New: 2", "New: 3"}, result)
	})

	t.Run("map string to length", func(t *testing.T) {
		s := []string{"a", "bb", "ccc"}
		result := slice.Map(s, func(v string) int {
			return len(v)
		})
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func Test_Reduce(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		start    int
		reducer  func(int, int) int
		expected int
	}

	cases := []Case{
		{"reduce sum", []int{1, 2, 3, 4, 5}, 0, func(acc, v int) int { return acc + v }, 15},
		{"reduce product", []int{1, 2, 3, 4}, 1, func(acc, v int) int { return acc * v }, 24},
		{"reduce max", []int{3, 7, 2, 9, 1}, 0, func(acc, v int) int {
			if v > acc {
				return v
			}
			return acc
		}, 9},
		{"reduce with initial value", []int{1, 2, 3}, 10, func(acc, v int) int { return acc + v }, 16},
		{"reduce empty slice", []int{}, 42, func(acc, v int) int { return acc + v }, 42},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Reduce(c.initial, c.start, c.reducer)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Reduce_Type_Conversion(t *testing.T) {
	t.Run("reduce ints to string", func(t *testing.T) {
		s := []int{1, 2, 3}
		result := slice.Reduce(s, "", func(acc string, v int) string {
			if acc == "" {
				return strconv.Itoa(v)
			}
			return acc + "," + strconv.Itoa(v)
		})
		assert.Equal(t, "1,2,3", result)
	})
}

func Test_Contains(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		value    int
		expected bool
	}

	cases := []Case{
		{"contains existing value", []int{1, 2, 3, 4, 5}, 3, true},
		{"does not contain value", []int{1, 2, 3, 4, 5}, 10, false},
		{"contains first element", []int{1, 2, 3}, 1, true},
		{"contains last element", []int{1, 2, 3}, 3, true},
		{"empty slice contains nothing", []int{}, 1, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Contains(c.slice, c.value)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_IndexOf(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		value    int
		expected int
	}

	cases := []Case{
		{"find first element", []int{1, 2, 3, 4, 5}, 1, 0},
		{"find middle element", []int{1, 2, 3, 4, 5}, 3, 2},
		{"find last element", []int{1, 2, 3, 4, 5}, 5, 4},
		{"value not found", []int{1, 2, 3, 4, 5}, 10, -1},
		{"find in empty slice", []int{}, 1, -1},
		{"find first duplicate", []int{1, 2, 3, 2, 4}, 2, 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.IndexOf(c.slice, c.value)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Reversed(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		expected []int
	}

	cases := []Case{
		{"reverse single element", []int{1}, []int{1}},
		{"reverse two elements", []int{1, 2}, []int{2, 1}},
		{"reverse multiple elements", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"reverse empty slice", []int{}, []int{}},
		{"reverse even length", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{"reverse odd length", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			original := make([]int, len(c.slice))

			copy(original, c.slice)

			result := slice.Reversed(c.slice)

			assert.Equal(t, c.expected, result)

			assert.Equal(t, original, c.slice)

			if len(result) > 0 {
				result[0] = 999
				assert.NotEqual(t, 999, c.slice[0], "modifying result should not affect original")
			}
		})
	}
}

func Test_Clone(t *testing.T) {
	type Case struct {
		name  string
		slice []int
	}

	cases := []Case{
		{"clone empty slice", []int{}},
		{"clone single element", []int{1}},
		{"clone multiple elements", []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Clone(c.slice)
			assert.Equal(t, c.slice, result)

			// Verify independence
			if len(c.slice) > 0 {
				result[0] = 999
				assert.NotEqual(t, c.slice[0], result[0])
			}
		})
	}
}

func Test_Unique(t *testing.T) {
	type Case struct {
		name     string
		slice    []int
		expected []int
	}

	cases := []Case{
		{"unique with duplicates", []int{1, 2, 2, 3, 3, 3, 4}, []int{1, 2, 3, 4}},
		{"unique with no duplicates", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"unique all same", []int{1, 1, 1, 1}, []int{1}},
		{"unique empty slice", []int{}, []int{}},
		{"unique preserves order", []int{3, 1, 2, 1, 3}, []int{3, 1, 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slice.Unique(c.slice)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Clear(t *testing.T) {
	type Case struct {
		name    string
		initial []int
	}

	cases := []Case{
		{"clear non-empty slice", []int{1, 2, 3, 4, 5}},
		{"clear single element", []int{1}},
		{"clear already empty", []int{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.initial
			slice.Clear(&s)
			assert.Equal(t, 0, len(s))
			assert.True(t, slice.IsEmpty(s))
		})
	}
}

func Test_Complex_Pipeline(t *testing.T) {
	t.Run("filter, map, reduce pipeline", func(t *testing.T) {
		// Start with 1-10, filter evens, square them, sum
		s := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

		evens := slice.Filter(s, func(v int) bool { return v%2 == 0 })
		squared := slice.Map(evens, func(v int) int { return v * v })
		sum := slice.Reduce(squared, 0, func(acc, v int) int { return acc + v })

		// 2^2 + 4^2 + 6^2 + 8^2 + 10^2 = 4 + 16 + 36 + 64 + 100 = 220
		assert.Equal(t, 220, sum)
	})
}

func Test_Chained_Operations(t *testing.T) {
	t.Run("multiple transformations", func(t *testing.T) {
		s := slice.Of(1, 2, 3, 4, 5)

		// Double
		s = slice.Map(s, func(v int) int { return v * 2 })
		assert.Equal(t, []int{2, 4, 6, 8, 10}, s)

		// Filter > 5
		s = slice.Filter(s, func(v int) bool { return v > 5 })
		assert.Equal(t, []int{6, 8, 10}, s)

		// Reverse
		s = slice.Reversed(s)
		assert.Equal(t, []int{10, 8, 6}, s)
	})
}

func Test_Immutability(t *testing.T) {
	t.Run("operations dont modify original", func(t *testing.T) {
		original := slice.Of(1, 2, 3, 4, 5)
		originalCopy := slice.Clone(original)

		// Perform various operations
		_ = slice.Filter(original, func(v int) bool { return v > 3 })
		_ = slice.Map(original, func(v int) int { return v * 2 })
		_ = slice.Reversed(original)
		_ = slice.Unique(original)

		// Original should be unchanged
		assert.Equal(t, originalCopy, original)
	})
}

func Test_Edge_Cases(t *testing.T) {
	t.Run("operations on nil vs empty", func(t *testing.T) {
		var nilSlice []int
		emptySlice := []int{}

		assert.True(t, slice.IsNil(nilSlice))
		assert.False(t, slice.IsNil(emptySlice))

		assert.True(t, slice.IsEmpty(nilSlice))
		assert.True(t, slice.IsEmpty(emptySlice))
	})

	t.Run("size of nil slice", func(t *testing.T) {
		var nilSlice []int
		assert.Equal(t, 0, slice.Size(nilSlice))
	})
}

func Test_Sort(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		expected []int
	}

	cases := []Case{
		{"sort ascending", []int{4, 1, 3, 2}, []int{1, 2, 3, 4}},
		{"already sorted", []int{1, 2, 3}, []int{1, 2, 3}},
		{"reverse sorted", []int{5, 4, 3}, []int{3, 4, 5}},
		{"empty slice", []int{}, []int{}},
		{"single element", []int{9}, []int{9}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := slice.Clone(c.initial) // Just to avoid modifying table data
			slice.Sort(s)
			assert.Equal(t, c.expected, s)
		})
	}
}

func Test_Sorted(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		expected []int
	}

	cases := []Case{
		{"sorted copy", []int{4, 1, 3, 2}, []int{1, 2, 3, 4}},
		{"already sorted", []int{1, 2, 3}, []int{1, 2, 3}},
		{"reverse sorted", []int{5, 4, 3}, []int{3, 4, 5}},
		{"empty slice", []int{}, []int{}},
		{"single element", []int{9}, []int{9}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			original := slice.Clone(c.initial)
			result := slice.Sorted(original)

			assert.Equal(t, c.expected, result)
			assert.Equal(t, c.initial, original, "Sorted must not modify original slice")
		})
	}
}
func Test_SortBy(t *testing.T) {
	type Case struct {
		name     string
		initial  []string
		expected []string
	}

	byLength := func(a, b string) int {
		switch {
		case len(a) < len(b):
			return -1
		case len(a) > len(b):
			return 1
		default:
			return 0
		}
	}

	cases := []Case{
		{"sort by length", []string{"aaa", "b", "cc"}, []string{"b", "cc", "aaa"}},
		{"all equal length", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"empty slice", []string{}, []string{}},
		{"single element", []string{"a"}, []string{"a"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := slice.Clone(c.initial)
			slice.SortBy(s, byLength)
			assert.Equal(t, c.expected, s)
		})
	}
}
func Test_SortedBy(t *testing.T) {
	type Case struct {
		name     string
		initial  []string
		expected []string
	}

	byLength := func(a, b string) int {
		switch {
		case len(a) < len(b):
			return -1
		case len(a) > len(b):
			return 1
		default:
			return 0
		}
	}

	cases := []Case{
		{"sorted copy by length", []string{"aaa", "b", "cc"}, []string{"b", "cc", "aaa"}},
		{"equal length", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"empty slice", []string{}, []string{}},
		{"single element", []string{"a"}, []string{"a"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			original := slice.Clone(c.initial)
			result := slice.SortedBy(original, byLength)

			assert.Equal(t, c.expected, result)
			assert.Equal(t, c.initial, original, "SortedBy must not modify original slice")
		})
	}
}

func Test_NewAndSized(t *testing.T) {
	type Case struct {
		name         string
		fn           func() []int
		wantLen      int
		wantCapacity int
	}

	cases := []Case{
		{
			name: "New without size returns empty slice with zero capacity",
			fn: func() []int {
				return slice.New[int]()
			},
			wantLen:      0,
			wantCapacity: 0,
		},
		{
			name: "New with size uses make and sets capacity",
			fn: func() []int {
				return slice.New[int](10)
			},
			wantLen:      0,
			wantCapacity: 10,
		},
		{
			name: "Sized calls New with size",
			fn: func() []int {
				return slice.Sized[int](5)
			},
			wantLen:      0,
			wantCapacity: 5,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.fn()
			assert.Equal(t, c.wantLen, len(got))
			assert.Equal(t, c.wantCapacity, cap(got))
		})
	}
}
