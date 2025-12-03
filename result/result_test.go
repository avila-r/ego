package result_test

import (
	"fmt"
	"testing"

	"github.com/avila-r/ego/result"
	"github.com/avila-r/ego/failure"
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
		{"of with value and error", 42, failure.New("test error"), false, true},
		{"of with zero value and no error", 0, nil, true, false},
		{"of with zero value and error", 0, failure.New("test error"), false, true},
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
		{"error with simple error", failure.New("simple error")},
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
	err := failure.New("failure error")
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
		{"value from error result", result.Error[int](failure.New("error")), nil},
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
		{"error from error result", result.Error[int](failure.New("test")), true},
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

func Test_OkOr(t *testing.T) {
	type Case struct {
		name        string
		result      result.Result[int]
		fallbackErr error
		isError     bool
		expectedErr string
	}

	cases := []Case{
		{
			name:        "ok or with success returns original",
			result:      result.Ok(42),
			fallbackErr: failure.New("fallback error"),
			isError:     false,
			expectedErr: "",
		},
		{
			name:        "ok or with error replaces with new error",
			result:      result.Error[int](failure.New("original error")),
			fallbackErr: failure.New("replacement error"),
			isError:     true,
			expectedErr: "replacement error",
		},
		{
			name:        "ok or with empty result uses fallback",
			result:      result.Result[int]{},
			fallbackErr: failure.New("empty result error"),
			isError:     true,
			expectedErr: "empty result error",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.result.OkOr(c.fallbackErr)
			assert.Equal(t, c.isError, r.IsError())

			if c.isError {
				assert.Equal(t, c.expectedErr, r.Error().Error())
			} else {
				assert.True(t, r.IsSuccess())
				assert.Equal(t, 42, r.Unwrap())
			}
		})
	}
}

func Test_OkOr_WithDifferentTypes(t *testing.T) {
	t.Run("ok or with string type", func(t *testing.T) {
		r := result.Ok("hello")
		result := r.OkOr(failure.New("should not use this"))

		assert.True(t, result.IsSuccess())
		assert.Equal(t, "hello", result.Unwrap())
	})

	t.Run("ok or with user struct", func(t *testing.T) {
		user := User{ID: 1, Name: "Alice"}
		r := result.Ok(user)
		result := r.OkOr(failure.New("user not found"))

		assert.True(t, result.IsSuccess())
		assert.Equal(t, user, result.Unwrap())
	})

	t.Run("ok or replaces error in user result", func(t *testing.T) {
		r := result.Error[User](failure.New("database connection failed"))
		result := r.OkOr(failure.New("service unavailable"))

		assert.True(t, result.IsError())
		assert.Equal(t, "service unavailable", result.Error().Error())
	})
}

func Test_OkOr_ErrorReplacement_Scenarios(t *testing.T) {
	t.Run("replace generic error with specific error", func(t *testing.T) {
		r := result.Error[int](failure.New("something went wrong"))
		result := r.OkOr(failure.New("database connection timeout"))

		assert.True(t, result.IsError())
		assert.Equal(t, "database connection timeout", result.Error().Error())
	})

	t.Run("provide user-friendly error message", func(t *testing.T) {
		r := result.Error[string](failure.New("sql: no rows in result set"))
		result := r.OkOr(failure.New("user not found"))

		assert.True(t, result.IsError())
		assert.Equal(t, "user not found", result.Error().Error())
	})

	t.Run("chain ok or with other operations", func(t *testing.T) {
		r := result.Error[int](failure.New("initial error")).
			OkOr(failure.New("replaced error")).
			OnFailure(func(err error) {
				assert.Equal(t, "replaced error", err.Error())
			})

		assert.True(t, r.IsError())
	})
}

func Test_Success_Method(t *testing.T) {
	type Case struct {
		name     string
		initial  result.Result[int]
		newValue int
	}

	cases := []Case{
		{"success on empty result", result.Result[int]{}, 42},
		{"success on error result", result.Error[int](failure.New("error")), 100},
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
		{"ok on error result", result.Error[int](failure.New("error")), 100},
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
		{"failure on empty result", result.Result[int]{}, failure.New("new error")},
		{"failure on success result", result.Ok(42), failure.New("override error")},
		{"failure on existing error", result.Error[int](failure.New("old")), failure.New("new")},
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
		{"error result", result.Error[int](failure.New("error")), true},
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
		{"error result", result.Error[int](failure.New("error")), false},
		{"empty result", result.Result[int]{}, false},
		{"result with value and error", result.Of(42, failure.New("error")), false},
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
		{"error result", result.Error[int](failure.New("error")), true},
		{"empty result", result.Result[int]{}, true},
		{"result with value and error", result.Of(42, failure.New("error")), true},
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
		r := result.Error[int](failure.New("error"))
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
		r := result.Error[int](failure.New("error"))
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
		r := result.Error[int](failure.New("error"))
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
		r := result.Error[int](failure.New("original error"))
		assert.PanicsWithValue(t, "custom panic message", func() {
			_ = r.Expect("custom panic message")
		})
	})

	t.Run("expect without message uses error message", func(t *testing.T) {
		r := result.Error[int](failure.New("error message"))
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
		r = r.Failure(failure.New("changed to error"))
		assert.True(t, r.IsError())
	})

	t.Run("chain failure to success", func(t *testing.T) {
		r := result.Error[int](failure.New("error"))
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
			result.Error[int](failure.New("error")),
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
		r := result.Error[int](failure.New("first error"))
		r = r.Failure(failure.New("second error"))
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
	if err := failure.New("database error"); err != nil {
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
	return r.Failure(failure.New("invalid user"))
}

func simulateDBQuery(success bool) result.Result[User] {
	if !success {
		return result.Error[User](failure.New("database connection failed"))
	}
	return result.Ok(User{ID: 1, Name: "Alice"})
}

func Test_OnSuccess(t *testing.T) {
	t.Run("onSuccess executes on success", func(t *testing.T) {
		executed := false
		r := result.Ok(42).OnSuccess(func(v int) {
			executed = true
			assert.Equal(t, 42, v)
		})

		assert.True(t, executed)
		assert.True(t, r.IsSuccess())
	})

	t.Run("onSuccess does not execute on error", func(t *testing.T) {
		executed := false
		r := result.Error[int](failure.New("error")).OnSuccess(func(v int) {
			executed = true
		})

		assert.False(t, executed)
		assert.True(t, r.IsError())
	})
}

func Test_OnFailure(t *testing.T) {
	t.Run("onFailure executes on error", func(t *testing.T) {
		executed := false
		expectedErr := failure.New("test error")

		r := result.Error[int](expectedErr).OnFailure(func(err error) {
			executed = true
			assert.Equal(t, expectedErr, err)
		})

		assert.True(t, executed)
		assert.True(t, r.IsError())
	})

	t.Run("onFailure does not execute on success", func(t *testing.T) {
		executed := false
		r := result.Ok(42).OnFailure(func(err error) {
			executed = true
		})

		assert.False(t, executed)
		assert.True(t, r.IsSuccess())
	})
}

func Test_Inspect(t *testing.T) {
	t.Run("inspect allows side effects", func(t *testing.T) {
		inspected := 0
		r := result.Ok(42).
			Inspect(func(v int) {
				inspected = v
			})

		assert.Equal(t, 42, inspected)
		assert.True(t, r.IsSuccess())
		assert.Equal(t, 42, r.Unwrap())
	})
}

func Test_InspectErr(t *testing.T) {
	t.Run("inspectErr allows error inspection", func(t *testing.T) {
		var inspectedErr error
		expectedErr := failure.New("test error")

		r := result.Error[int](expectedErr).
			InspectErr(func(err error) {
				inspectedErr = err
			})

		assert.Equal(t, expectedErr, inspectedErr)
		assert.True(t, r.IsError())
	})
}

func Test_Or(t *testing.T) {
	type Case struct {
		name     string
		first    result.Result[int]
		second   result.Result[int]
		expected int
		isError  bool
	}

	cases := []Case{
		{
			name:     "or returns first if success",
			first:    result.Ok(42),
			second:   result.Ok(100),
			expected: 42,
			isError:  false,
		},
		{
			name:     "or returns second if first is error",
			first:    result.Error[int](failure.New("error")),
			second:   result.Ok(100),
			expected: 100,
			isError:  false,
		},
		{
			name:     "or returns second error if both error",
			first:    result.Error[int](failure.New("error1")),
			second:   result.Error[int](failure.New("error2")),
			expected: 0,
			isError:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.first.Or(c.second)
			if c.isError {
				assert.True(t, r.IsError())
			} else {
				assert.True(t, r.IsSuccess())
				assert.Equal(t, c.expected, r.Unwrap())
			}
		})
	}
}

func Test_OrElse(t *testing.T) {
	t.Run("orElse returns original on success", func(t *testing.T) {
		r := result.Ok(42).OrElse(func(err error) result.Result[int] {
			return result.Ok(100)
		})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 42, r.Unwrap())
	})

	t.Run("orElse calls function on error", func(t *testing.T) {
		r := result.Error[int](failure.New("error")).OrElse(func(err error) result.Result[int] {
			return result.Ok(100)
		})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 100, r.Unwrap())
	})
}

func Test_And(t *testing.T) {
	type Case struct {
		name     string
		first    result.Result[int]
		second   result.Result[int]
		expected int
		isError  bool
	}

	cases := []Case{
		{
			name:     "and returns second if first is success",
			first:    result.Ok(42),
			second:   result.Ok(100),
			expected: 100,
			isError:  false,
		},
		{
			name:     "and returns first error",
			first:    result.Error[int](failure.New("error1")),
			second:   result.Ok(100),
			expected: 0,
			isError:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.first.And(c.second)
			if c.isError {
				assert.True(t, r.IsError())
			} else {
				assert.True(t, r.IsSuccess())
				assert.Equal(t, c.expected, r.Unwrap())
			}
		})
	}
}

func Test_UnwrapOr(t *testing.T) {
	type Case struct {
		name         string
		result       result.Result[int]
		defaultValue int
		expected     int
	}

	cases := []Case{
		{
			name:         "unwrapOr returns value on success",
			result:       result.Ok(42),
			defaultValue: 999,
			expected:     42,
		},
		{
			name:         "unwrapOr returns default on error",
			result:       result.Error[int](failure.New("error")),
			defaultValue: 999,
			expected:     999,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := c.result.UnwrapOr(c.defaultValue)
			assert.Equal(t, c.expected, v)
		})
	}
}

func Test_UnwrapOrElse(t *testing.T) {
	t.Run("unwrapOrElse returns value on success", func(t *testing.T) {
		v := result.Ok(42).UnwrapOrElse(func(err error) int {
			return 999
		})
		assert.Equal(t, 42, v)
	})

	t.Run("unwrapOrElse calls function on error", func(t *testing.T) {
		v := result.Error[int](failure.New("error")).UnwrapOrElse(func(err error) int {
			return 999
		})
		assert.Equal(t, 999, v)
	})
}

func Test_UnwrapOrDefault(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		expected int
	}

	cases := []Case{
		{
			name:     "unwrapOrDefault returns value on success",
			result:   result.Ok(42),
			expected: 42,
		},
		{
			name:     "unwrapOrDefault returns zero value on error",
			result:   result.Error[int](failure.New("error")),
			expected: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := c.result.UnwrapOrDefault()
			assert.Equal(t, c.expected, v)
		})
	}
}

func Test_Contains(t *testing.T) {
	equals := func(a, b int) bool { return a == b }

	type Case struct {
		name     string
		result   result.Result[int]
		value    int
		expected bool
	}

	cases := []Case{
		{
			name:     "contains finds matching value",
			result:   result.Ok(42),
			value:    42,
			expected: true,
		},
		{
			name:     "contains returns false for different value",
			result:   result.Ok(42),
			value:    100,
			expected: false,
		},
		{
			name:     "contains returns false on error",
			result:   result.Error[int](failure.New("error")),
			value:    42,
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.result.Contains(c.value, equals))
		})
	}
}

func Test_ContainsErr(t *testing.T) {
	t.Run("containsErr finds matching error", func(t *testing.T) {
		targetErr := failure.New("specific error")
		r := result.Error[int](targetErr)

		contains := r.ContainsErr(func(err error) bool {
			return err == targetErr
		})

		assert.True(t, contains)
	})

	t.Run("containsErr returns false on success", func(t *testing.T) {
		r := result.Ok(42)

		contains := r.ContainsErr(func(err error) bool {
			return true
		})

		assert.False(t, contains)
	})
}

func Test_Match(t *testing.T) {
	t.Run("match executes success branch", func(t *testing.T) {
		successCalled := false
		failureCalled := false

		result.Ok(42).Match(
			func(v int) {
				successCalled = true
				assert.Equal(t, 42, v)
			},
			func(err error) {
				failureCalled = true
			},
		)

		assert.True(t, successCalled)
		assert.False(t, failureCalled)
	})

	t.Run("match executes failure branch", func(t *testing.T) {
		successCalled := false
		failureCalled := false
		expectedErr := failure.New("test error")

		result.Error[int](expectedErr).Match(
			func(v int) {
				successCalled = true
			},
			func(err error) {
				failureCalled = true
				assert.Equal(t, expectedErr, err)
			},
		)

		assert.False(t, successCalled)
		assert.True(t, failureCalled)
	})
}

func Test_MatchReturn(t *testing.T) {
	t.Run("matchReturn returns from success branch", func(t *testing.T) {
		v := result.MatchReturn(
			result.Ok(42),
			func(v int) string {
				return "success"
			},
			func(err error) string {
				return "failure"
			},
		)

		assert.Equal(t, "success", v)
	})

	t.Run("matchReturn returns from failure branch", func(t *testing.T) {
		v := result.MatchReturn(
			result.Error[int](failure.New("error")),
			func(v int) string {
				return "success"
			},
			func(err error) string {
				return "failure"
			},
		)

		assert.Equal(t, "failure", v)
	})
}

func Test_Map(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		mapFn    func(int) int
		expected int
		isError  bool
	}

	cases := []Case{
		{
			name:     "map success value",
			result:   result.Ok(5),
			mapFn:    func(v int) int { return v * 2 },
			expected: 10,
			isError:  false,
		},
		{
			name:     "map with zero value",
			result:   result.Ok(0),
			mapFn:    func(v int) int { return v + 10 },
			expected: 10,
			isError:  false,
		},
		{
			name:     "map on error result does not apply function",
			result:   result.Error[int](failure.New("error")),
			mapFn:    func(v int) int { return v * 2 },
			expected: 0,
			isError:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.result.Map(c.mapFn)
			assert.Equal(t, c.isError, r.IsError())
			if !c.isError {
				assert.Equal(t, c.expected, r.Unwrap())
			}
		})
	}
}

func Test_MapOr(t *testing.T) {
	type Case struct {
		name         string
		result       result.Result[int]
		defaultValue int
		mapFn        func(int) int
		expected     int
	}

	cases := []Case{
		{
			name:         "map or with success",
			result:       result.Ok(5),
			defaultValue: 100,
			mapFn:        func(v int) int { return v * 2 },
			expected:     10,
		},
		{
			name:         "map or with error returns default",
			result:       result.Error[int](failure.New("error")),
			defaultValue: 100,
			mapFn:        func(v int) int { return v * 2 },
			expected:     100,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := c.result.MapOr(c.defaultValue, c.mapFn)
			assert.Equal(t, c.expected, v)
		})
	}
}

func Test_MapOrElse(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		fallback func(error) int
		mapFn    func(int) int
		expected int
	}

	cases := []Case{
		{
			name:     "map or else with success",
			result:   result.Ok(5),
			fallback: func(e error) int { return -1 },
			mapFn:    func(v int) int { return v * 2 },
			expected: 10,
		},
		{
			name:     "map or else with error calls fallback",
			result:   result.Error[int](failure.New("error")),
			fallback: func(e error) int { return 999 },
			mapFn:    func(v int) int { return v * 2 },
			expected: 999,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := c.result.MapOrElse(c.fallback, c.mapFn)
			assert.Equal(t, c.expected, v)
		})
	}
}

func Test_FlatMap(t *testing.T) {
	type Case struct {
		name     string
		result   result.Result[int]
		flatMap  func(int) result.Result[int]
		expected int
		isError  bool
	}

	cases := []Case{
		{
			name:   "flatmap success to success",
			result: result.Ok(5),
			flatMap: func(v int) result.Result[int] {
				return result.Ok(v * 2)
			},
			expected: 10,
			isError:  false,
		},
		{
			name:   "flatmap success to error",
			result: result.Ok(5),
			flatMap: func(v int) result.Result[int] {
				return result.Error[int](failure.New("flatmap error"))
			},
			expected: 0,
			isError:  true,
		},
		{
			name:   "flatmap error does not call function",
			result: result.Error[int](failure.New("initial error")),
			flatMap: func(v int) result.Result[int] {
				return result.Ok(v * 2)
			},
			expected: 0,
			isError:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.result.FlatMap(c.flatMap)
			assert.Equal(t, c.isError, r.IsError())
			if !c.isError {
				assert.Equal(t, c.expected, r.Unwrap())
			}
		})
	}
}

func Test_AndThen(t *testing.T) {
	t.Run("andthen chains successfully", func(t *testing.T) {
		r := result.Ok(5).
			AndThen(func(v int) result.Result[int] {
				return result.Ok(v * 2)
			}).
			AndThen(func(v int) result.Result[int] {
				return result.Ok(v + 10)
			})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 20, r.Unwrap())
	})

	t.Run("andthen stops at first error", func(t *testing.T) {
		r := result.Ok(5).
			AndThen(func(v int) result.Result[int] {
				return result.Error[int](failure.New("chain broken"))
			}).
			AndThen(func(v int) result.Result[int] {
				return result.Ok(v + 10)
			})

		assert.True(t, r.IsError())
		assert.Equal(t, "chain broken", r.Error().Error())
	})
}

func Test_Bind(t *testing.T) {
	t.Run("bind is alias for flatmap", func(t *testing.T) {
		r1 := result.Ok(5).Bind(func(v int) result.Result[int] {
			return result.Ok(v * 2)
		})

		r2 := result.Ok(5).FlatMap(func(v int) result.Result[int] {
			return result.Ok(v * 2)
		})

		assert.Equal(t, r1.Unwrap(), r2.Unwrap())
	})
}

func Test_Filter(t *testing.T) {
	type Case struct {
		name      string
		result    result.Result[int]
		predicate func(int) bool
		err       error
		isError   bool
	}

	cases := []Case{
		{
			name:      "filter passes",
			result:    result.Ok(10),
			predicate: func(v int) bool { return v > 5 },
			err:       failure.New("filter failed"),
			isError:   false,
		},
		{
			name:      "filter fails",
			result:    result.Ok(3),
			predicate: func(v int) bool { return v > 5 },
			err:       failure.New("value too small"),
			isError:   true,
		},
		{
			name:      "filter on error returns error",
			result:    result.Error[int](failure.New("initial error")),
			predicate: func(v int) bool { return true },
			err:       failure.New("filter failed"),
			isError:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := c.result.Filter(c.predicate, c.err)
			assert.Equal(t, c.isError, r.IsError())
		})
	}
}

func Test_Try(t *testing.T) {
	t.Run("try with successful function", func(t *testing.T) {
		r := result.Try(func() int {
			return 42
		})
		assert.True(t, r.IsSuccess())
		assert.Equal(t, 42, r.Unwrap())
	})

	t.Run("try catches panic", func(t *testing.T) {
		r := result.Try(func() int {
			panic(failure.New("panic error"))
		})
		assert.True(t, r.IsError())
		assert.Error(t, r.Error())
	})

	t.Run("try catches string panic", func(t *testing.T) {
		r := result.Try(func() int {
			panic("string panic")
		})
		assert.True(t, r.IsError())
	})
}

func Test_TryWith(t *testing.T) {
	type Case struct {
		name    string
		fn      func() (int, error)
		isError bool
		value   int
	}

	cases := []Case{
		{
			name: "try with success",
			fn: func() (int, error) {
				return 42, nil
			},
			isError: false,
			value:   42,
		},
		{
			name: "try with error",
			fn: func() (int, error) {
				return 0, failure.New("function error")
			},
			isError: true,
			value:   0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := result.TryWith(c.fn)
			assert.Equal(t, c.isError, r.IsError())
			if !c.isError {
				assert.Equal(t, c.value, r.Unwrap())
			}
		})
	}
}

func Test_Chain(t *testing.T) {
	t.Run("chain multiple operations", func(t *testing.T) {
		initial := result.Ok(5)

		r := result.Chain(
			initial,
			func(r result.Result[int]) result.Result[int] {
				return r.Map(func(v int) int { return v * 2 })
			},
			func(r result.Result[int]) result.Result[int] {
				return r.Map(func(v int) int { return v + 10 })
			},
		)

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 20, r.Unwrap())
	})

	t.Run("chain stops at first error", func(t *testing.T) {
		initial := result.Ok(5)

		r := result.Chain(
			initial,
			func(r result.Result[int]) result.Result[int] {
				return result.Error[int](failure.New("chain error"))
			},
			func(r result.Result[int]) result.Result[int] {
				return r.Map(func(v int) int { return v + 10 })
			},
		)

		assert.True(t, r.IsError())
		assert.Equal(t, "chain error", r.Error().Error())
	})
}

func Test_Flatten(t *testing.T) {
	t.Run("flatten nested success", func(t *testing.T) {
		nested := result.Ok(result.Ok(42))
		flattened := result.Flatten(nested)

		assert.True(t, flattened.IsSuccess())
		assert.Equal(t, 42, flattened.Unwrap())
	})

	t.Run("flatten outer error", func(t *testing.T) {
		nested := result.Error[result.Result[int]](failure.New("outer error"))
		flattened := result.Flatten(nested)

		assert.True(t, flattened.IsError())
		assert.Equal(t, "outer error", flattened.Error().Error())
	})

	t.Run("flatten inner error", func(t *testing.T) {
		nested := result.Ok(result.Error[int](failure.New("inner error")))
		flattened := result.Flatten(nested)

		assert.True(t, flattened.IsError())
		assert.Equal(t, "inner error", flattened.Error().Error())
	})
}

func Test_Pipeline_Scenario(t *testing.T) {
	t.Run("validation and transformation pipeline", func(t *testing.T) {
		validateAge := func(age int) result.Result[int] {
			if age < 0 {
				return result.Error[int](failure.New("age cannot be negative"))
			}
			if age > 150 {
				return result.Error[int](failure.New("age too high"))
			}
			return result.Ok(age)
		}

		r := result.
			TryWith(func() (int, error) {
				return 25, nil
			}).
			AndThen(validateAge).
			Map(func(age int) int { return age + 1 }).
			Filter(func(age int) bool {
				return age >= 18
			}, failure.New("must be adult"))

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 26, r.Unwrap())
	})

	t.Run("database operation pipeline", func(t *testing.T) {
		fetchUser := func(id int) result.Result[User] {
			if id <= 0 {
				return result.Error[User](failure.New("invalid id"))
			}
			return result.Ok(User{ID: id, Name: "Alice"})
		}

		validateUser := func(u User) result.Result[User] {
			if u.Name == "" {
				return result.Error[User](failure.New("name is required"))
			}
			return result.Ok(u)
		}

		r := result.FlatMap(result.FlatMap(result.Ok(1), fetchUser), validateUser)

		assert.True(t, r.IsSuccess())
		assert.Equal(t, "Alice", r.Unwrap().Name)
	})
}

func Test_Complex_Chaining(t *testing.T) {
	t.Run("complex map and flatmap chain", func(t *testing.T) {
		parseInput := func(s string) result.Result[int] {
			if s == "" {
				return result.Error[int](failure.New("empty input"))
			}
			return result.Ok(len(s))
		}

		r := result.FlatMap(result.Ok("hello"), parseInput).
			Map(func(v int) int {
				return v * 2
			}).
			Filter(func(v int) bool {
				return v > 5
			}, failure.New("too small"))

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 10, r.Unwrap())
	})
}

func Test_Recover(t *testing.T) {
	t.Run("recover from error", func(t *testing.T) {
		r := result.Error[int](failure.New("error")).Recover(func(err error) int {
			return 100
		})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 100, r.Unwrap())
	})

	t.Run("recover on success does nothing", func(t *testing.T) {
		r := result.Ok(42).Recover(func(err error) int {
			return 100
		})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 42, r.Unwrap())
	})
}

func Test_RecoverWith(t *testing.T) {
	t.Run("recover with from error", func(t *testing.T) {
		r := result.Error[int](failure.New("error")).RecoverWith(func(err error) result.Result[int] {
			return result.Ok(100)
		})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 100, r.Unwrap())
	})

	t.Run("recover with on success does nothing", func(t *testing.T) {
		r := result.Ok(42).RecoverWith(func(err error) result.Result[int] {
			return result.Ok(100)
		})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 42, r.Unwrap())
	})
}

func Test_Transpose(t *testing.T) {
	t.Run("transpose success", func(t *testing.T) {
		r := result.Ok(42)
		v, err := r.Transpose()

		assert.NoError(t, err)
		assert.Equal(t, 42, *v)
	})

	t.Run("transpose error", func(t *testing.T) {
		expectedErr := failure.New("test error")
		r := result.Error[int](expectedErr)
		v, err := r.Transpose()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, v)
	})
}

func Test_FromPointer(t *testing.T) {
	t.Run("from pointer with value", func(t *testing.T) {
		v := 42
		r := result.FromPointer(&v, nil)

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 42, r.Unwrap())
	})

	t.Run("from pointer with error", func(t *testing.T) {
		v := 42
		expectedErr := failure.New("error")
		r := result.FromPointer(&v, expectedErr)

		assert.True(t, r.IsError())
		assert.Equal(t, expectedErr, r.Error())
	})

	t.Run("from pointer with nil value", func(t *testing.T) {
		r := result.FromPointer[int](nil, nil)
		assert.True(t, r.IsError())
	})
}

func Test_ToPointer(t *testing.T) {
	t.Run("to pointer from success", func(t *testing.T) {
		r := result.Ok(42)
		v, err := r.ToPointer()

		assert.NoError(t, err)
		assert.Equal(t, 42, *v)
	})

	t.Run("to pointer from error", func(t *testing.T) {
		r := result.Error[int](failure.New("error"))
		v, err := r.ToPointer()

		assert.Error(t, err)
		assert.Nil(t, v)
	})
}

func Test_Combine(t *testing.T) {
	t.Run("combine two successes", func(t *testing.T) {
		r1 := result.Ok(42)
		r2 := result.Ok("hello")

		combined := result.Combine(r1, r2)

		assert.True(t, combined.IsSuccess())
		assert.Equal(t, 42, combined.Unwrap().First)
		assert.Equal(t, "hello", combined.Unwrap().Second)
	})

	t.Run("combine with first error", func(t *testing.T) {
		r1 := result.Error[int](failure.New("first error"))
		r2 := result.Ok("hello")

		combined := result.Combine(r1, r2)

		assert.True(t, combined.IsError())
		assert.Equal(t, "first error", combined.Error().Error())
	})

	t.Run("combine with second error", func(t *testing.T) {
		r1 := result.Ok(42)
		r2 := result.Error[string](failure.New("second error"))

		combined := result.Combine(r1, r2)

		assert.True(t, combined.IsError())
		assert.Equal(t, "second error", combined.Error().Error())
	})
}

func Test_MapErr(t *testing.T) {
	t.Run("map err transforms error", func(t *testing.T) {
		r := result.Error[int](failure.New("original")).MapErr(func(err error) error {
			return failure.New("transformed: %s", err.Error())
		})

		assert.True(t, r.IsError())
		assert.Equal(t, "transformed: original", r.Error().Error())
	})

	t.Run("map err on success does nothing", func(t *testing.T) {
		r := result.Ok(42).MapErr(func(err error) error {
			return failure.New("should not be called")
		})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 42, r.Unwrap())
	})
}

func Test_ExpectErr(t *testing.T) {
	t.Run("expect err on error returns error", func(t *testing.T) {
		expectedErr := failure.New("test error")
		r := result.Error[int](expectedErr)
		err := r.ExpectErr("should not panic")

		assert.Equal(t, expectedErr, err)
	})

	t.Run("expect err on success panics", func(t *testing.T) {
		r := result.Ok(42)
		assert.PanicsWithValue(t, "expected an error", func() {
			_ = r.ExpectErr("expected an error")
		})
	})
}

func Test_UnwrapErr(t *testing.T) {
	t.Run("unwrap err on error returns error", func(t *testing.T) {
		expectedErr := failure.New("test error")
		r := result.Error[int](expectedErr)
		err := r.UnwrapErr()

		assert.Equal(t, expectedErr, err)
	})

	t.Run("unwrap err on success panics", func(t *testing.T) {
		r := result.Ok(42)
		assert.Panics(t, func() {
			_ = r.UnwrapErr()
		})
	})
}

func Test_Iter(t *testing.T) {
	t.Run("iter on success returns slice with value", func(t *testing.T) {
		r := result.Ok(42)
		slice := r.Iter()

		assert.Len(t, slice, 1)
		assert.Equal(t, 42, slice[0])
	})

	t.Run("iter on error returns empty slice", func(t *testing.T) {
		r := result.Error[int](failure.New("error"))
		slice := r.Iter()

		assert.Len(t, slice, 0)
	})
}

func Test_Collect(t *testing.T) {
	t.Run("collect all successes", func(t *testing.T) {
		results := []result.Result[int]{
			result.Ok(1),
			result.Ok(2),
			result.Ok(3),
		}

		r := result.Collect(results)

		assert.True(t, r.IsSuccess())
		assert.Equal(t, []int{1, 2, 3}, r.Unwrap())
	})

	t.Run("collect with error stops at first error", func(t *testing.T) {
		results := []result.Result[int]{
			result.Ok(1),
			result.Error[int](failure.New("error")),
			result.Ok(3),
		}

		r := result.Collect(results)

		assert.True(t, r.IsError())
	})
}

func Test_Partition(t *testing.T) {
	t.Run("partition mixed results", func(t *testing.T) {
		results := []result.Result[int]{
			result.Ok(1),
			result.Error[int](failure.New("error1")),
			result.Ok(2),
			result.Error[int](failure.New("error2")),
			result.Ok(3),
		}

		successes, failures := result.Partition(results)

		assert.Equal(t, []int{1, 2, 3}, successes)
		assert.Len(t, failures, 2)
	})
}

func Test_Tap(t *testing.T) {
	t.Run("tap executes side effect and returns result", func(t *testing.T) {
		tapped := false
		r := result.Ok(42).Tap(func(r result.Result[int]) {
			tapped = true
			assert.True(t, r.IsSuccess())
		})

		assert.True(t, tapped)
		assert.True(t, r.IsSuccess())
		assert.Equal(t, 42, r.Unwrap())
	})
}

func Test_RealWorld_Scenarios(t *testing.T) {
	t.Run("user registration pipeline", func(t *testing.T) {
		validateEmail := func(email string) result.Result[string] {
			if email == "" {
				return result.Error[string](failure.New("email required"))
			}
			return result.Ok(email)
		}

		validatePassword := func(pwd string) result.Result[string] {
			if len(pwd) < 8 {
				return result.Error[string](failure.New("password too short"))
			}
			return result.Ok(pwd)
		}

		r := result.Ok("user@example.com").
			FlatMap(validateEmail).
			FlatMap(validatePassword).
			OnSuccess(func(email string) {
				// Log success
			}).
			Map(func(email string) string {
				return email + ":validated"
			})

		assert.True(t, r.IsSuccess())
	})

	t.Run("error recovery chain", func(t *testing.T) {
		r := result.Error[int](failure.New("connection failed")).
			RecoverWith(func(err error) result.Result[int] {
				// Try cache
				return result.Ok(999)
			}).
			OnSuccess(func(v int) {
				// Log cache hit
			})

		assert.True(t, r.IsSuccess())
		assert.Equal(t, 999, r.Unwrap())
	})
}
