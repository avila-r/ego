package optional_test

import (
	"testing"

	"github.com/avila-r/ego/optional"
	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	type Case struct {
		name     string
		value    int
		expected int
	}

	cases := []Case{
		{"create with zero value", 0, 0},
		{"create with positive value", 42, 42},
		{"create with negative value", -10, -10},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			opt := optional.Of(c.value)
			assert.True(t, opt.IsPresent())
			assert.False(t, opt.IsEmpty())
			val, ok := opt.Get()
			assert.True(t, ok)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Of_WithString(t *testing.T) {
	opt := optional.Of("hello")
	assert.True(t, opt.IsPresent())
	val, ok := opt.Get()
	assert.True(t, ok)
	assert.Equal(t, "hello", val)
}

func Test_Of_WithStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "John", Age: 30}
	opt := optional.Of(person)
	assert.True(t, opt.IsPresent())
	val, ok := opt.Get()
	assert.True(t, ok)
	assert.Equal(t, person, val)
}

func Test_Empty(t *testing.T) {
	type Case struct {
		name string
	}

	cases := []Case{
		{"create empty int optional"},
		{"create empty string optional"},
		{"create empty struct optional"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			opt := optional.Empty[int]()
			assert.False(t, opt.IsPresent())
			assert.True(t, opt.IsEmpty())
			_, ok := opt.Get()
			assert.False(t, ok)
		})
	}
}

func Test_IsPresent(t *testing.T) {
	type Case struct {
		name     string
		opt      optional.Optional[int]
		expected bool
	}

	cases := []Case{
		{"present value", optional.Of(42), true},
		{"empty optional", optional.Empty[int](), false},
		{"zero value is present", optional.Of(0), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.opt.IsPresent())
		})
	}
}

func Test_IsEmpty(t *testing.T) {
	type Case struct {
		name     string
		opt      optional.Optional[int]
		expected bool
	}

	cases := []Case{
		{"empty optional", optional.Empty[int](), true},
		{"present value", optional.Of(42), false},
		{"zero value is not empty", optional.Of(0), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.opt.IsEmpty())
		})
	}
}

func Test_Join(t *testing.T) {
	type Case struct {
		name     string
		opt      optional.Optional[int]
		expected int
	}

	cases := []Case{
		{"join present value", optional.Of(42), 42},
		{"join zero value", optional.Of(0), 0},
		{"join negative value", optional.Of(-5), -5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val := c.opt.Join()
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Join_Panics(t *testing.T) {
	opt := optional.Empty[int]()
	assert.Panics(t, func() {
		opt.Join()
	})
}

func Test_Get(t *testing.T) {
	type Case struct {
		name        string
		opt         optional.Optional[int]
		expectedVal int
		expectedOk  bool
	}

	cases := []Case{
		{"get present value", optional.Of(42), 42, true},
		{"get zero value", optional.Of(0), 0, true},
		{"get from empty", optional.Empty[int](), 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val, ok := c.opt.Get()
			assert.Equal(t, c.expectedOk, ok)
			assert.Equal(t, c.expectedVal, val)
		})
	}
}

func Test_Clear(t *testing.T) {
	type Case struct {
		name string
		opt  optional.Optional[int]
	}

	cases := []Case{
		{"clear present value", optional.Of(42)},
		{"clear already empty", optional.Empty[int]()},
		{"clear zero value", optional.Of(0)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.opt.Clear()
			assert.True(t, c.opt.IsEmpty())
			assert.False(t, c.opt.IsPresent())
			_, ok := c.opt.Get()
			assert.False(t, ok)
		})
	}
}

func Test_GetOrDefault(t *testing.T) {
	type Case struct {
		name     string
		opt      optional.Optional[int]
		fallback int
		expected int
	}

	cases := []Case{
		{"present value returns value", optional.Of(42), 99, 42},
		{"empty returns fallback", optional.Empty[int](), 99, 99},
		{"zero value returns zero", optional.Of(0), 99, 0},
		{"empty with zero fallback", optional.Empty[int](), 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val := c.opt.GetOrDefault(c.fallback)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Set(t *testing.T) {
	type Case struct {
		name     string
		initial  optional.Optional[int]
		newValue int
		expected int
	}

	cases := []Case{
		{"set on empty optional", optional.Empty[int](), 42, 42},
		{"set on present optional", optional.Of(10), 20, 20},
		{"set zero value", optional.Of(10), 0, 0},
		{"overwrite with negative", optional.Of(5), -5, -5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.initial.Set(c.newValue)
			assert.True(t, c.initial.IsPresent())
			val, ok := c.initial.Get()
			assert.True(t, ok)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Take(t *testing.T) {
	type Case struct {
		name        string
		opt         optional.Optional[int]
		expectedVal int
		expectError bool
	}

	cases := []Case{
		{"take present value", optional.Of(42), 42, false},
		{"take zero value", optional.Of(0), 0, false},
		{"take from empty", optional.Empty[int](), 0, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val, err := c.opt.Take()
			if c.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, val)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, val)
				assert.Equal(t, c.expectedVal, *val)
			}
		})
	}
}

func Test_Optional_ComplexWorkflow(t *testing.T) {
	// Start with empty
	opt := optional.Empty[string]()
	assert.True(t, opt.IsEmpty())

	// Set a value
	opt.Set("hello")
	assert.True(t, opt.IsPresent())
	assert.Equal(t, "hello", opt.Join())

	// Get the value
	val, ok := opt.Get()
	assert.True(t, ok)
	assert.Equal(t, "hello", val)

	// Take the value
	ptr, err := opt.Take()
	assert.Nil(t, err)
	assert.Equal(t, "hello", *ptr)

	// Clear it
	opt.Clear()
	assert.True(t, opt.IsEmpty())

	// GetOrDefault after clear
	val = opt.GetOrDefault("default")
	assert.Equal(t, "default", val)
}

func Test_Optional_Pointer_Independence(t *testing.T) {
	// Test that modifying returned pointer doesn't affect optional
	opt := optional.Of(42)

	ptr, err := opt.Take()
	assert.Nil(t, err)
	assert.Equal(t, 42, *ptr)

	// Original optional should still have value
	assert.True(t, opt.IsPresent())
	assert.Equal(t, 42, opt.Join())
}

func Test_Optional_MultipleSet(t *testing.T) {
	opt := optional.Empty[int]()

	opt.Set(1)
	assert.Equal(t, 1, opt.Join())

	opt.Set(2)
	assert.Equal(t, 2, opt.Join())

	opt.Set(3)
	assert.Equal(t, 3, opt.Join())
}

func Test_Optional_StringType(t *testing.T) {
	opt := optional.Of("test")
	assert.True(t, opt.IsPresent())
	assert.Equal(t, "test", opt.Join())

	opt.Clear()
	assert.True(t, opt.IsEmpty())
	assert.Equal(t, "fallback", opt.GetOrDefault("fallback"))
}

func Test_Optional_BoolType(t *testing.T) {
	optTrue := optional.Of(true)
	assert.True(t, optTrue.IsPresent())
	assert.True(t, optTrue.Join())

	optFalse := optional.Of(false)
	assert.True(t, optFalse.IsPresent())
	assert.False(t, optFalse.Join())
}
