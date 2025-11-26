package box_test

import (
	"fmt"
	"testing"

	"github.com/avila-r/ego/box"
	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	type Case struct {
		name     string
		value    int
		expected int
	}

	cases := []Case{
		{"create box with positive value", 42, 42},
		{"create box with zero", 0, 0},
		{"create box with negative value", -10, -10},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := box.Of(c.value)
			assert.True(t, b.IsPresent())
			assert.False(t, b.IsEmpty())
			assert.Equal(t, c.expected, *b.Get())
		})
	}
}

func Test_Empty(t *testing.T) {
	b := box.Empty[int]()
	assert.False(t, b.IsPresent())
	assert.True(t, b.IsEmpty())
	assert.Nil(t, b.Get())
}

func Test_Get(t *testing.T) {
	type Case struct {
		name     string
		box      box.Box[int]
		expected int
	}

	value := 42
	cases := []Case{
		{"get from present box", box.Of(42), value},
		{"get from empty box", box.Empty[int](), -1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.box.Get()
			if c.expected == -1 {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, c.expected, *result)
			}
		})
	}
}

func Test_GetOrDefault(t *testing.T) {
	type Case struct {
		name         string
		box          box.Box[int]
		defaultValue int
		expected     int
	}

	cases := []Case{
		{"present box returns value", box.Of(42), 100, 42},
		{"empty box returns default", box.Empty[int](), 100, 100},
		{"empty box with zero default", box.Empty[int](), 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.box.GetOrDefault(c.defaultValue)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_With(t *testing.T) {
	type Case struct {
		name     string
		initial  box.Box[int]
		newValue int
		expected int
	}

	cases := []Case{
		{"set value on empty box", box.Empty[int](), 42, 42},
		{"replace value on present box", box.Of(10), 20, 20},
		{"set zero value", box.Empty[int](), 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.initial.With(c.newValue)
			assert.True(t, result.IsPresent())
			assert.Equal(t, c.expected, *result.Get())
		})
	}
}

func Test_Peek(t *testing.T) {
	t.Run("peek on present box executes action", func(t *testing.T) {
		called := false
		var capturedValue int

		b := box.Of(42)
		result := b.Peek(func(v int) {
			called = true
			capturedValue = v
		})

		assert.True(t, called)
		assert.Equal(t, 42, capturedValue)
		assert.Equal(t, b, result) // Returns same box
	})

	t.Run("peek on empty box does not execute action", func(t *testing.T) {
		called := false

		b := box.Empty[int]()
		result := b.Peek(func(v int) {
			called = true
		})

		assert.False(t, called)
		assert.Equal(t, b, result)
	})
}

func Test_Then(t *testing.T) {
	type Case struct {
		name     string
		initial  box.Box[int]
		mapper   func(int) int
		expected int
	}

	value := 84
	cases := []Case{
		{"then on present box applies mapper", box.Of(42), func(v int) int { return v * 2 }, value},
		{"then on empty box does nothing", box.Empty[int](), func(v int) int { return v * 2 }, -1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.initial.Then(c.mapper)
			if c.expected == -1 {
				assert.True(t, result.IsEmpty())
			} else {
				assert.Equal(t, c.expected, *result.Get())
			}
		})
	}
}

func Test_ThenSupplier(t *testing.T) {
	type Case struct {
		name     string
		initial  box.Box[int]
		supplier func() int
		expected int
	}

	cases := []Case{
		{"supplier on empty box", box.Empty[int](), func() int { return 100 }, 100},
		{"supplier on present box replaces value", box.Of(42), func() int { return 200 }, 200},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.initial.ThenSupplier(c.supplier)
			assert.True(t, result.IsPresent())
			assert.Equal(t, c.expected, *result.Get())
		})
	}
}

func Test_ThenConsumer(t *testing.T) {
	t.Run("consumer on present box executes action", func(t *testing.T) {
		called := false
		var capturedValue int

		b := box.Of(42)
		result := b.ThenConsumer(func(v int) {
			called = true
			capturedValue = v
		})

		assert.True(t, called)
		assert.Equal(t, 42, capturedValue)
		assert.Equal(t, b, result)
	})

	t.Run("consumer on empty box does not execute action", func(t *testing.T) {
		called := false

		b := box.Empty[int]()
		result := b.ThenConsumer(func(v int) {
			called = true
		})

		assert.False(t, called)
		assert.Equal(t, b, result)
	})
}

func Test_Filter(t *testing.T) {
	type Case struct {
		name          string
		initial       box.Box[int]
		predicate     func(int) bool
		shouldBeEmpty bool
	}

	cases := []Case{
		{"filter passes on present box", box.Of(42), func(v int) bool { return v > 40 }, false},
		{"filter fails on present box", box.Of(42), func(v int) bool { return v < 40 }, true},
		{"filter on empty box returns empty", box.Empty[int](), func(v int) bool { return true }, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.initial.Filter(c.predicate)
			if c.shouldBeEmpty {
				assert.True(t, result.IsEmpty())
			} else {
				assert.True(t, result.IsPresent())
			}
		})
	}
}

func Test_Map(t *testing.T) {
	type Case struct {
		name     string
		initial  box.Box[int]
		mapper   func(int) string
		expected *string
	}

	strValue := "42"
	cases := []Case{
		{"map on present box", box.Of(42), func(v int) string { return fmt.Sprintf("%d", v) }, &strValue},
		{"map on empty box", box.Empty[int](), func(v int) string { return fmt.Sprintf("%d", v) }, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Map(c.initial, c.mapper)
			if c.expected == nil {
				assert.True(t, result.IsEmpty())
			} else {
				assert.Equal(t, *c.expected, *result.Get())
			}
		})
	}
}

func Test_FlatMap(t *testing.T) {
	type Case struct {
		name        string
		initial     box.Box[int]
		mapper      func(int) box.Box[string]
		shouldEmpty bool
		expected    string
	}

	cases := []Case{
		{
			"flatmap on present box to present box",
			box.Of(42),
			func(v int) box.Box[string] { return box.Of(fmt.Sprintf("%d", v)) },
			false,
			"42",
		},
		{
			"flatmap on present box to empty box",
			box.Of(42),
			func(v int) box.Box[string] { return box.Empty[string]() },
			true,
			"",
		},
		{
			"flatmap on empty box",
			box.Empty[int](),
			func(v int) box.Box[string] { return box.Of(fmt.Sprintf("%d", v)) },
			true,
			"",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.FlatMap(c.initial, c.mapper)
			if c.shouldEmpty {
				assert.True(t, result.IsEmpty())
			} else {
				assert.Equal(t, c.expected, *result.Get())
			}
		})
	}
}

func Test_Deflated(t *testing.T) {
	type Case struct {
		name    string
		initial box.Box[int]
	}

	cases := []Case{
		{"deflate present box", box.Of(42)},
		{"deflate empty box", box.Empty[int]()},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.initial.Deflated()
			assert.True(t, result.IsEmpty())
			assert.False(t, result.IsPresent())
			assert.Nil(t, result.Get())
		})
	}
}

func Test_Copy(t *testing.T) {
	type Case struct {
		name          string
		initial       box.Box[int]
		shouldBeEmpty bool
	}

	cases := []Case{
		{"copy present box", box.Of(42), false},
		{"copy empty box", box.Empty[int](), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.initial.Copy()
			if c.shouldBeEmpty {
				assert.True(t, result.IsEmpty())
			} else {
				assert.True(t, result.IsPresent())
				assert.Equal(t, *c.initial.Get(), *result.Get())
				// Verify it's a different box instance (pointer comparison)
				assert.NotSame(t, &c.initial, &result)

				// Verify independence: modifying copy doesn't affect original
				result.Set(999)
				assert.Equal(t, 42, *c.initial.Get())
				assert.Equal(t, 999, *result.Get())
			}
		})
	}
}

func Test_Copy_Independence(t *testing.T) {
	original := box.Of(42)
	copied := original.Copy()

	// Modify the copy
	copied.Set(100)

	// Original should be unchanged
	assert.Equal(t, 42, *original.Get())
	assert.Equal(t, 100, *copied.Get())
}

func Test_Set(t *testing.T) {
	type Case struct {
		name     string
		initial  box.Box[int]
		newValue int
	}

	cases := []Case{
		{"set on empty box", box.Empty[int](), 42},
		{"set on present box", box.Of(10), 20},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.initial.Set(c.newValue)
			assert.True(t, c.initial.IsPresent())
			assert.Equal(t, c.newValue, *c.initial.Get())
		})
	}
}

func Test_Deflate(t *testing.T) {
	type Case struct {
		name    string
		initial box.Box[int]
	}

	cases := []Case{
		{"deflate present box", box.Of(42)},
		{"deflate empty box", box.Empty[int]()},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.initial.Deflate()
			assert.True(t, c.initial.IsEmpty())
			assert.Nil(t, c.initial.Get())
		})
	}
}

func Test_IfPresent(t *testing.T) {
	t.Run("if present on present box executes action", func(t *testing.T) {
		called := false
		var capturedValue int

		b := box.Of(42)
		b.IfPresent(func(v int) {
			called = true
			capturedValue = v
		})

		assert.True(t, called)
		assert.Equal(t, 42, capturedValue)
	})

	t.Run("if present on empty box does not execute action", func(t *testing.T) {
		called := false

		b := box.Empty[int]()
		b.IfPresent(func(v int) {
			called = true
		})

		assert.False(t, called)
	})
}

func Test_String(t *testing.T) {
	type Case struct {
		name     string
		box      box.Box[int]
		expected string
	}

	cases := []Case{
		{"string representation of present box", box.Of(42), "Box(42)"},
		{"string representation of empty box", box.Empty[int](), "Box(<empty>)"},
		{"string representation of zero value", box.Of(0), "Box(0)"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.box.String()
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Equals(t *testing.T) {
	type Case struct {
		name     string
		box1     box.Box[int]
		box2     box.Box[int]
		expected bool
	}

	cases := []Case{
		{"same instance", box.Of(42), box.Of(42), true},
		{"different present boxes", box.Of(42), box.Of(10), false},
		{"both empty boxes", box.Empty[int](), box.Empty[int](), true},
		{"present and empty box", box.Of(42), box.Empty[int](), false},
		{"empty and present box", box.Empty[int](), box.Of(42), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.box1.Equals(&c.box2)
			assert.Equal(t, c.expected, result)
		})
	}

	t.Run("equals with self", func(t *testing.T) {
		b := box.Of(42)
		assert.True(t, b.Equals(&b))
	})
}

func Test_Increment(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value5 := 5
	value0 := 0
	valueNeg := -3

	cases := []Case{
		{"increment positive value", value5, 6},
		{"increment zero", value0, 1},
		{"increment negative value", valueNeg, -2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Increment(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Decrement(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value5 := 5
	value0 := 0
	valueNeg := -3

	cases := []Case{
		{"decrement positive value", value5, 4},
		{"decrement zero", value0, -1},
		{"decrement negative value", valueNeg, -4},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Decrement(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Square(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value5 := 5
	value0 := 0
	valueNeg := -3

	cases := []Case{
		{"square positive value", value5, 25},
		{"square zero", value0, 0},
		{"square negative value", valueNeg, 9},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Square(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Cube(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value3 := 3
	value0 := 0
	valueNeg := -2

	cases := []Case{
		{"cube positive value", value3, 27},
		{"cube zero", value0, 0},
		{"cube negative value", valueNeg, -8},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Cube(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Twice(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value5 := 5
	value0 := 0
	valueNeg := -3

	cases := []Case{
		{"twice positive value", value5, 10},
		{"twice zero", value0, 0},
		{"twice negative value", valueNeg, -6},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Twice(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Halve(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value10 := 10
	value5 := 5
	value0 := 0
	valueNeg := -6

	cases := []Case{
		{"halve even value", value10, 5},
		{"halve odd value truncates", value5, 2},
		{"halve zero", value0, 0},
		{"halve negative value", valueNeg, -3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Halve(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Negate(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value5 := 5
	value0 := 0
	valueNeg := -3

	cases := []Case{
		{"negate positive value", value5, -5},
		{"negate zero", value0, 0},
		{"negate negative value", valueNeg, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Negate(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Abs(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value5 := 5
	value0 := 0
	valueNeg := -3

	cases := []Case{
		{"abs of positive value", value5, 5},
		{"abs of zero", value0, 0},
		{"abs of negative value", valueNeg, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Abs(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Identity(t *testing.T) {
	type Case struct {
		name     string
		input    int
		expected int
	}

	value5 := 5
	value0 := 0
	valueNeg := -3

	cases := []Case{
		{"identity of positive value", value5, 5},
		{"identity of zero", value0, 0},
		{"identity of negative value", valueNeg, -3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Identity(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Modulo(t *testing.T) {
	type Case struct {
		name     string
		input    int
		modulus  int
		expected int
	}

	value10 := 10
	value7 := 7
	value0 := 0

	cases := []Case{
		{"modulo even division", value10, 5, 0},
		{"modulo with remainder", value10, 3, 1},
		{"modulo larger than value", value7, 10, 7},
		{"modulo of zero", value0, 5, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Modulo(c.modulus)(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Clamp(t *testing.T) {
	type Case struct {
		name     string
		input    int
		min      int
		max      int
		expected int
	}

	value5 := 5
	value100 := 100
	valueNeg := -10
	value50 := 50

	cases := []Case{
		{"clamp within range", value50, 0, 100, 50},
		{"clamp below min", value5, 10, 100, 10},
		{"clamp above max", value100, 10, 50, 50},
		{"clamp negative below min", valueNeg, 0, 100, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := box.Clamp(c.min, c.max)(c.input)
			assert.Equal(t, c.expected, result)
		})
	}
}

func Test_Box_With_Increment(t *testing.T) {
	b := box.Of(5)
	result := b.Then(func(v int) int {
		return box.Increment(v)
	})
	assert.Equal(t, 6, *result.Get())
}

func Test_Box_With_Decrement(t *testing.T) {
	b := box.Of(10)
	result := b.Then(func(v int) int {
		return box.Decrement(v)
	})
	assert.Equal(t, 9, *result.Get())
}

func Test_Box_With_Square(t *testing.T) {
	b := box.Of(5)
	result := b.Then(func(v int) int {
		return box.Square(v)
	})
	assert.Equal(t, 25, *result.Get())
}

func Test_Box_With_Cube(t *testing.T) {
	b := box.Of(3)
	result := b.Then(func(v int) int {
		return box.Cube(v)
	})
	assert.Equal(t, 27, *result.Get())
}

func Test_Box_With_Twice(t *testing.T) {
	b := box.Of(7)
	result := b.Then(func(v int) int {
		return box.Twice(v)
	})
	assert.Equal(t, 14, *result.Get())
}

func Test_Box_With_Halve(t *testing.T) {
	b := box.Of(20)
	result := b.Then(func(v int) int {
		return box.Halve(v)
	})
	assert.Equal(t, 10, *result.Get())
}

func Test_Box_With_Negate(t *testing.T) {
	b := box.Of(15)
	result := b.Then(func(v int) int {
		return box.Negate(v)
	})
	assert.Equal(t, -15, *result.Get())
}

func Test_Box_With_Abs(t *testing.T) {
	b := box.Of(-25)
	result := b.Then(func(v int) int {
		return box.Abs(v)
	})
	assert.Equal(t, 25, *result.Get())
}

func Test_Box_With_Modulo(t *testing.T) {
	b := box.Of(17)
	result := b.Then(func(v int) int {
		return box.Modulo(5)(v)
	})
	assert.Equal(t, 2, *result.Get())
}

func Test_Box_With_Clamp(t *testing.T) {
	type Case struct {
		name     string
		value    int
		min      int
		max      int
		expected int
	}

	cases := []Case{
		{"clamp value below range", 5, 10, 20, 10},
		{"clamp value above range", 30, 10, 20, 20},
		{"clamp value within range", 15, 10, 20, 15},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := box.Of(c.value)
			result := b.Then(func(v int) int {
				return box.Clamp(c.min, c.max)(v)
			})
			assert.Equal(t, c.expected, *result.Get())
		})
	}
}

func Test_Box_Chain_Multiple_Operations(t *testing.T) {
	// Start with 5, square it (25), then halve it (12), then increment (13)
	b := box.Of(5)
	result := b.
		Then(func(v int) int { return box.Square(v) }).
		Then(func(v int) int { return box.Halve(v) }).
		Then(func(v int) int { return box.Increment(v) })

	assert.Equal(t, 13, *result.Get())
}

func Test_Box_Empty_With_Toolkit(t *testing.T) {
	b := box.Empty[int]()
	result := b.Then(func(v int) int {
		return box.Square(v)
	})
	assert.True(t, result.IsEmpty())
}

func Test_Box_Filter_With_Toolkit(t *testing.T) {
	// Filter out odd numbers, then square
	b := box.Of(4)
	result := b.
		Filter(func(v int) bool { return v%2 == 0 }).
		Then(func(v int) int { return box.Square(v) })

	assert.Equal(t, 16, *result.Get())

	// Filter should fail for odd number
	b2 := box.Of(5)
	result2 := b2.
		Filter(func(v int) bool { return v%2 == 0 }).
		Then(func(v int) int { return box.Square(v) })

	assert.True(t, result2.IsEmpty())
}

func Test_Box_Peek_With_Toolkit(t *testing.T) {
	var sideEffect int
	b := box.Of(10)

	result := b.
		Peek(func(v int) {
			sideEffect = box.Twice(v)
		}).
		Then(func(v int) int { return box.Increment(v) })

	assert.Equal(t, 20, sideEffect)    // Side effect from peek
	assert.Equal(t, 11, *result.Get()) // Original value incremented
}

func Test_Box_Complex_Workflow(t *testing.T) {
	// Complex workflow: start with 10, double it, clamp to 15, square, modulo 100
	b := box.Of(10)
	result := b.
		Then(func(v int) int { return box.Twice(v) }).      // 20
		Then(func(v int) int { return box.Clamp(0, 15)(v) }). // 15
		Then(func(v int) int { return box.Square(v) }).       // 225
		Then(func(v int) int { return box.Modulo(100)(v) })   // 25

	assert.Equal(t, 25, *result.Get())
}
