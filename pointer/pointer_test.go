package pointer_test

import (
	"testing"

	"github.com/avila-r/ego/pointer"
	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	type Case[T any] struct {
		name  string
		input T
	}

	cases := []Case[int]{
		{"int value", 10},
		{"zero int", 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ptr := pointer.Of(c.input)
			assert.NotNil(t, ptr)
			assert.Equal(t, c.input, *ptr)
		})
	}

	// separate type tests (string, struct, bool etc.)
	t.Run("string", func(t *testing.T) {
		ptr := pointer.Of("abc")
		assert.NotNil(t, ptr)
		assert.Equal(t, "abc", *ptr)
	})

	t.Run("bool", func(t *testing.T) {
		ptr := pointer.Of(true)
		assert.NotNil(t, ptr)
		assert.Equal(t, true, *ptr)
	})

	t.Run("struct", func(t *testing.T) {
		type S struct {
			X int
		}
		in := S{X: 5}
		ptr := pointer.Of(in)
		assert.NotNil(t, ptr)
		assert.Equal(t, 5, ptr.X)
	})
}
