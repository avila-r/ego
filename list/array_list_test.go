package list_test

import (
	"testing"

	"github.com/avila-r/ego/list"
	"github.com/avila-r/ego/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		toAdd    []int
		expected []int
	}

	cases := []Case{
		{"add single element", []int{}, []int{1}, []int{1}},
		{"add multiple elements", []int{1}, []int{2, 3}, []int{1, 2, 3}},
		{"add to non-empty list", []int{10, 20}, []int{30}, []int{10, 20, 30}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.NewArrayList[int]()
			l.Add(c.initial...)
			l.Add(c.toAdd...)
			assert.Equal(t, c.expected, l.Items())
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
		{"valid index", []int{1, 2, 3}, 1, 2, true},
		{"first element", []int{10, 20, 30}, 0, 10, true},
		{"last element", []int{10, 20, 30}, 2, 30, true},
		{"out of range positive", []int{1, 2}, 5, 0, false},
		{"out of range negative", []int{1, 2}, -1, 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.EmptyArrayList[int]()
			l.Add(c.init...)
			val, ok := l.Get(c.index)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Set(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		index    int
		newVal   int
		expected []int
		success  bool
	}

	cases := []Case{
		{"set middle element", []int{1, 2, 3}, 1, 99, []int{1, 99, 3}, true},
		{"set first element", []int{5, 6}, 0, 10, []int{10, 6}, true},
		{"set last element", []int{5, 6, 7}, 2, 70, []int{5, 6, 70}, true},
		{"set out of range", []int{1, 2}, 5, 100, []int{1, 2}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.EmptyArrayList[int]()
			l.Add(c.init...)
			ok := l.Set(c.index, c.newVal)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, l.Items())
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
		{"remove invalid index", []int{1, 2}, 5, []int{1, 2}, false},
		{"remove from empty list", []int{}, 0, []int{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.EmptyArrayList[int]()
			l.Add(c.init...)
			ok := l.Remove(c.index)
			assert.Equal(t, c.success, ok)
			assert.Equal(t, c.expected, l.Items())
		})
	}
}

func Test_Contains(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		value    int
		expected bool
	}

	cases := []Case{
		{"contains value", []int{1, 2, 3}, 2, true},
		{"does not contain value", []int{1, 2, 3}, 4, false},
		{"empty list", []int{}, 10, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.EmptyArrayList[int]()
			l.Add(c.init...)
			assert.Equal(t, c.expected, l.Contains(c.value))
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
		{"empty list", []int{}, 0},
		{"three elements", []int{1, 2, 3}, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.EmptyArrayList[int]()
			l.Add(c.init...)
			assert.Equal(t, c.expected, l.Size())
		})
	}
}

func Test_Clear(t *testing.T) {
	l := list.EmptyArrayList[int]()
	l.Add(1, 2, 3)
	l.Clear()
	assert.Equal(t, 0, l.Size())
	assert.False(t, l.Contains(1))
	assert.False(t, l.Contains(2))
	assert.False(t, l.Contains(3))
}

func Test_Items(t *testing.T) {
	type Case struct {
		name     string
		init     []int
		expected []int
	}

	cases := []Case{
		{"return all elements", []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty list", []int{}, []int{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.EmptyArrayList[int]()
			l.Add(c.init...)
			assert.Equal(t, c.expected, l.Items())
		})
	}
}

func Test_Elements(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		expected []int
	}

	cases := []Case{
		{"elements from non-empty list", []int{1, 2, 3}, []int{1, 2, 3}},
		{"elements from empty list", []int{}, []int{}},
		{"elements single", []int{9}, []int{9}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.Empty[int]()
			l.Add(c.initial...)
			assert.Equal(t, c.expected, l.Elements())
		})
	}
}

func Test_ForEach(t *testing.T) {
	type Case struct {
		name     string
		input    []int
		expected int
	}

	cases := []Case{
		{"sum of elements", []int{1, 2, 3}, 6},
		{"sum of empty list", []int{}, 0},
		{"sum single", []int{10}, 10},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.EmptyArrayList[int]()
			l.Add(c.input...)

			sum := 0
			l.ForEach(func(v int) { sum += v })

			assert.Equal(t, c.expected, sum)
		})
	}
}

func Test_IsEmpty(t *testing.T) {
	type Case struct {
		name     string
		initial  []int
		expected bool
	}

	cases := []Case{
		{"empty list", []int{}, true},
		{"non-empty list", []int{1}, false},
		{"multiple elements", []int{1, 2}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.EmptyArrayList[int]()
			l.Add(c.initial...)
			assert.Equal(t, c.expected, l.IsEmpty())
		})
	}
}

func Test_Of(t *testing.T) {
	type Case struct {
		name     string
		input    []int
		expected []int
	}

	cases := []Case{
		{"of empty", []int{}, []int{}},
		{"of single", []int{5}, []int{5}},
		{"of multiple", []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.Of(c.input...)
			assert.Equal(t, c.expected, l.Elements())
		})
	}
}

func Test_Iterator(t *testing.T) {
	type Case struct {
		name     string
		input    []int
		expected []int
	}

	cases := []Case{
		{"empty list", []int{}, []int{}},
		{"single element", []int{1}, []int{1}},
		{"multiple elements", []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.NewArrayList[int]()

			l.Add(c.input...)

			it := l.Iterator()

			result := slice.Empty[int]()
			for it.HasNext() {
				result = append(result, it.Next())
			}

			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Stream(t *testing.T) {
	type Case struct {
		name     string
		input    []int
		expected []int
	}

	cases := []Case{
		{"empty list", []int{}, []int{}},
		{"single element", []int{10}, []int{10}},
		{"multiple elements", []int{1, 2, 3}, []int{1, 2, 3}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := list.NewArrayList[int]()
			l.Add(c.input...)

			s := l.Stream()
			out := s.ToSlice() // change if needed

			assert.Equal(t, c.expected, out)
		})
	}
}
