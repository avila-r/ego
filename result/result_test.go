package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/avila-r/ego/result"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int
	Name string
}

func Test_Of(t *testing.T) {
	type Case struct {
		name      string
		value     int
		err       error
		isSuccess bool
		isError   bool
	}

	cases := []Case{
		{"of with value and no error", 42, nil, true, false},
		{"of with value and error", 42, errors.New("test error"), false, true},
		{"of with zero value and no error", 0, nil, true, false},
		{"of with zero value and error", 0, errors.New("test error"), false, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := result.Of(c.value, c.err)
			assert.Equal(t, c.isSuccess, r.IsSuccess())
			assert.Equal(t, c.isError, r.IsError())
			if c.err == nil {
				assert.NoError(t, r.Error())
			} else {
				assert.Equal(t, c.err, r.Error())
			}
		})
	}
}

func Test_Ok(t *testing.T) {
	type Case struct {
		name  string
		value int
	}

	cases := []Case{
		{"ok with positive value", 42},
		{"ok with zero", 0},
		{"ok with negative value", -10},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := result.Ok(c.value)
			assert.True(t, r.IsSuccess())
			assert.False(t, r.IsError())
			assert.NoError(t, r.Error())
			assert.Equal(t, c.value, *r.Value())
		})
	}
}

func Test_Error(t *testing.T) {
	type Case struct {
		name string
		err  error
	}

	cases := []Case{
		{"error with simple error", errors.New("simple error")},
		{"error with formatted error", fmt.Errorf("formatted: %s", "error")},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := result.Error[int](c.err)
			assert.False(t, r.IsSuccess())
			assert.True(t, r.IsError())
			assert.Equal(t, c.err, r.Error())
			assert.Nil(t, r.Value())
		})
	}
}

func Test_Failure(t *testing.T) {
	err := errors.New("failure error")
	r := result.Failure[int](err)

	assert.False(t, r.IsSuccess())
	assert.True(t, r.IsError())
	assert.Equal(t, err, r.Error())
	assert.Nil(t, r.Value())
}

func Test_Value(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		expected *int
	}

	value := 42
	cases := []Case{
		{"value from success result", result.Ok(42), &value},
		{"value from error result", result.Error[int](errors.New("error")), nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := c.result.Value()
			if c.expected == nil {
				assert.Nil(t, v)
			} else {
				assert.Equal(t, *c.expected, *v)
			}
		})
	}
}

func Test_Error_Method(t *testing.T) {
	type Case struct {
		name        string
		result      result.Result[int]
		shouldError bool
	}

	cases := []Case{
		{"error from success result", result.Ok(42), false},
		{"error from error result", result.Error[int](errors.New("test")), true},
		{"error from empty result", result.Result[int]{}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.result.Error()
			if c.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Success_Method(t *testing.T) {
	type Case struct {
		name     string
		initial  result.Result[int]
		newValue int
	}

	cases := []Case{
		{"success on empty result", result.Result[int]{}, 42},
		{"success on error result", result.Error[int](errors.New("error")), 100},
		{"success on existing success", result.Ok(10), 20},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.initial.Success(c.newValue)
			assert.True(t, r.IsSuccess())
			assert.False(t, r.IsError())
			assert.Equal(t, c.newValue, *r.Value())
		})
	}
}

func Test_Ok_Method(t *testing.T) {
	type Case struct {
		name    string
		initial result.Result[int]
		value   int
	}

	cases := []Case{
		{"ok on empty result", result.Result[int]{}, 42},
		{"ok on error result", result.Error[int](errors.New("error")), 100},
		{"ok on existing success", result.Ok(10), 20},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.initial.Ok(c.value)
			assert.True(t, r.IsSuccess())
			assert.False(t, r.IsError())
			assert.Equal(t, c.value, *r.Value())
		})
	}
}

func Test_Failure_Method(t *testing.T) {
	type Case struct {
		name    string
		initial result.Result[int]
		err     error
	}
	cases := []Case{
		{"failure on empty result", result.Result[int]{}, errors.New("new error")},
		{"failure on success result", result.Ok(42), errors.New("override error")},
		{"failure on existing error", result.Error[int](errors.New("old")), errors.New("new")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.initial.Failure(c.err)
			assert.False(t, r.IsSuccess())
			assert.True(t, r.IsError())
			assert.Equal(t, c.err, r.Error())
		})
	}
}

func Test_IsEmpty(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		expected bool
	}

	cases := []Case{
		{"empty result", result.Result[int]{}, true},
		{"success result", result.Ok(42), false},
		{"error result", result.Error[int](errors.New("error")), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.result.IsEmpty())
		})
	}
}

func Test_IsSuccess(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		expected bool
	}

	cases := []Case{
		{"success result", result.Ok(42), true},
		{"error result", result.Error[int](errors.New("error")), false},
		{"empty result", result.Result[int]{}, false},
		{"result with value and error", result.Of(42, errors.New("error")), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.result.IsSuccess())
		})
	}
}

func Test_IsError(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		expected bool
	}

	cases := []Case{
		{"success result", result.Ok(42), false},
		{"error result", result.Error[int](errors.New("error")), true},
		{"empty result", result.Result[int]{}, true},
		{"result with value and error", result.Of(42, errors.New("error")), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.result.IsError())
		})
	}
}

func Test_Unwrap(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		expected int
	}

	cases := []Case{
		{"unwrap success", result.Ok(42), 42},
		{"unwrap zero value", result.Ok(0), 0},
		{"unwrap negative value", result.Ok(-10), -10},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.result.Unwrap())
		})
	}
}

func Test_Unwrap_Panics(t *testing.T) {
	t.Run("unwrap panics on error result", func(t *testing.T) {
		r := result.Error[int](errors.New("error"))
		assert.Panics(t, func() {
			_ = r.Unwrap()
		})
	})

	t.Run("unwrap panics on empty result", func(t *testing.T) {
		r := result.Result[int]{}
		assert.Panics(t, func() {
			_ = r.Unwrap()
		})
	})
}

func Test_Take(t *testing.T) {
	t.Run("take from success result", func(t *testing.T) {
		r := result.Ok(42)
		v, err := r.Take()
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, 42, *v)
	})

	t.Run("take from empty result", func(t *testing.T) {
		r := result.Result[int]{}
		v, err := r.Take()
		assert.Error(t, err)
		assert.Nil(t, v)
	})

	t.Run("take from error result", func(t *testing.T) {
		r := result.Error[int](errors.New("error"))
		v, err := r.Take()
		assert.Error(t, err)
		assert.Nil(t, v)
	})
}

func Test_Join(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		expected int
	}

	cases := []Case{
		{"join success", result.Ok(42), 42},
		{"join zero value", result.Ok(0), 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.result.Join())
		})
	}
}

func Test_Join_Panics(t *testing.T) {
	t.Run("join panics on error result", func(t *testing.T) {
		r := result.Error[int](errors.New("error"))
		assert.Panics(t, func() {
			_ = r.Join()
		})
	})

	t.Run("join panics on empty result", func(t *testing.T) {
		r := result.Result[int]{}
		assert.Panics(t, func() {
			_ = r.Join()
		})
	})
}

func Test_Expect(t *testing.T) {
	t.Run("expect with success returns value", func(t *testing.T) {
		r := result.Ok(42)
		assert.Equal(t, 42, r.Expect("should not panic"))
	})

	t.Run("expect with custom message panics", func(t *testing.T) {
		r := result.Error[int](errors.New("original error"))
		assert.PanicsWithValue(t, "custom panic message", func() {
			_ = r.Expect("custom panic message")
		})
	})

	t.Run("expect without message uses error message", func(t *testing.T) {
		r := result.Error[int](errors.New("error message"))
		assert.PanicsWithValue(t, "error message", func() {
			_ = r.Expect()
		})
	})

	t.Run("expect on empty result panics", func(t *testing.T) {
		r := result.Result[int]{}
		assert.Panics(t, func() {
			_ = r.Expect("empty result")
		})
	})
}

func Test_NamedReturn_UseCase(t *testing.T) {
	t.Run("named return with success", func(t *testing.T) {
		r := createUserSuccess(User{ID: 1, Name: "Alice"})
		assert.True(t, r.IsSuccess())
		assert.Equal(t, 1, r.Unwrap())
	})

	t.Run("named return with failure", func(t *testing.T) {
		r := createUserFailure(User{ID: 1, Name: "Alice"})
		assert.True(t, r.IsError())
		assert.Error(t, r.Error())
	})

	t.Run("named return forgotten attribution", func(t *testing.T) {
		r := createUserForgotten(User{ID: 1, Name: "Alice"})
		assert.True(t, r.IsEmpty())
		assert.True(t, r.IsError())
		assert.Error(t, r.Error())
	})

	t.Run("named return early success", func(t *testing.T) {
		r := createUserEarlyReturn(User{ID: 5, Name: "Bob"})
		assert.True(t, r.IsSuccess())
		assert.Equal(t, 5, r.Unwrap())
	})
}

func Test_Chaining_Operations(t *testing.T) {
	t.Run("chain success to failure", func(t *testing.T) {
		r := result.Ok(42)
		r = r.Failure(errors.New("changed to error"))
		assert.True(t, r.IsError())
	})

	t.Run("chain failure to success", func(t *testing.T) {
		r := result.Error[int](errors.New("error"))
		r = r.Success(100)
		assert.True(t, r.IsSuccess())
		assert.Equal(t, 100, r.Unwrap())
	})

	t.Run("chain multiple successes", func(t *testing.T) {
		r := result.Ok(10)
		r = r.Success(20)
		r = r.Success(30)
		assert.True(t, r.IsSuccess())
		assert.Equal(t, 30, r.Unwrap())
	})
}

func Test_Complex_Scenarios(t *testing.T) {
	t.Run("database operation simulation", func(t *testing.T) {
		// Simulate successful DB operation
		r := simulateDBQuery(true)
		assert.True(t, r.IsSuccess())
		user := r.Unwrap()
		assert.Equal(t, "Alice", user.Name)
	})

	t.Run("database operation failure", func(t *testing.T) {
		// Simulate failed DB operation
		r := simulateDBQuery(false)
		assert.True(t, r.IsError())
		assert.Error(t, r.Error())
	})

	t.Run("result in pipeline", func(t *testing.T) {
		results := []result.Result[int]{
			result.Ok(1),
			result.Ok(2),
			result.Error[int](errors.New("error")),
			result.Ok(4),
		}

		sum := 0
		errorCount := 0

		for _, r := range results {
			if r.IsSuccess() {
				sum += r.Unwrap()
			} else {
				errorCount++
			}
		}

		assert.Equal(t, 7, sum)
		assert.Equal(t, 1, errorCount)
	})
}

func Test_Edge_Cases(t *testing.T) {
	t.Run("result with struct type", func(t *testing.T) {
		user := User{ID: 1, Name: "Alice"}
		r := result.Ok(user)
		assert.True(t, r.IsSuccess())
		assert.Equal(t, user, r.Unwrap())
	})

	t.Run("result with pointer type", func(t *testing.T) {
		user := &User{ID: 1, Name: "Alice"}
		r := result.Ok(user)
		assert.True(t, r.IsSuccess())
		assert.Equal(t, user, r.Unwrap())
	})

	t.Run("result with string type", func(t *testing.T) {
		r := result.Ok("hello")
		assert.True(t, r.IsSuccess())
		assert.Equal(t, "hello", r.Unwrap())
	})

	t.Run("result with empty string", func(t *testing.T) {
		r := result.Ok("")
		assert.True(t, r.IsSuccess())
		assert.Equal(t, "", r.Unwrap())
	})

	t.Run("multiple errors", func(t *testing.T) {
		r := result.Error[int](errors.New("first error"))
		r = r.Failure(errors.New("second error"))
		assert.Equal(t, "second error", r.Error().Error())
	})
}

func Test_Zero_Values(t *testing.T) {
	t.Run("zero value int is valid success", func(t *testing.T) {
		r := result.Ok(0)
		assert.True(t, r.IsSuccess())
		assert.Equal(t, 0, r.Unwrap())
	})

	t.Run("zero value string is valid success", func(t *testing.T) {
		r := result.Ok("")
		assert.True(t, r.IsSuccess())
		assert.Equal(t, "", r.Unwrap())
	})

	t.Run("zero value bool is valid success", func(t *testing.T) {
		r := result.Ok(false)
		assert.True(t, r.IsSuccess())
		assert.Equal(t, false, r.Unwrap())
	})
}

func createUserSuccess(user User) (r result.Result[int]) {
	return r.Success(user.ID)
}

func createUserFailure(user User) (r result.Result[int]) {
	if err := errors.New("database error"); err != nil {
		return r.Failure(err)
	}
	return r.Success(user.ID)
}

func createUserForgotten(user User) (r result.Result[int]) {
	amazing := false
	if amazing {
		return r.Success(user.ID)
	}
	return // Forgotten to set Success or Failure
}

func createUserEarlyReturn(user User) (r result.Result[int]) {
	if user.ID > 0 {
		return r.Success(user.ID)
	}
	return r.Failure(errors.New("invalid user"))
}

func simulateDBQuery(success bool) result.Result[User] {
	if !success {
		return result.Error[User](errors.New("database connection failed"))
	}
	return result.Ok(User{ID: 1, Name: "Alice"})
}
