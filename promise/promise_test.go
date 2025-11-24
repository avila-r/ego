package promise_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/avila-r/ego/promise"
	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	type Case struct {
		name     string
		supplier func() (int, error)
		expected int
		hasError bool
	}

	cases := []Case{
		{
			name: "successful completion",
			supplier: func() (int, error) {
				return 42, nil
			},
			expected: 42,
			hasError: false,
		},
		{
			name: "completion with error",
			supplier: func() (int, error) {
				return 0, errors.New("test error")
			},
			expected: 0,
			hasError: true,
		},
		{
			name: "delayed completion",
			supplier: func() (int, error) {
				time.Sleep(50 * time.Millisecond)
				return 100, nil
			},
			expected: 100,
			hasError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier)
			val, err := p.Get()
			assert.Equal(t, c.expected, val)
			assert.Equal(t, c.hasError, err != nil)
		})
	}
}

func Test_Run(t *testing.T) {
	type Case struct {
		name     string
		executed bool
	}

	cases := []Case{
		{"executes runnable", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			executed := c.executed
			p := promise.Run(func() {
				executed = true
			})
			p.Get()
			assert.True(t, executed)
		})
	}
}

func Test_Completed(t *testing.T) {
	type Case struct {
		name     string
		value    string
		expected string
	}

	cases := []Case{
		{"already completed", "hello", "hello"},
		{"empty string", "", ""},
		{"number as string", "123", "123"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Completed(c.value)
			val, err := p.Get()
			assert.Equal(t, c.expected, val)
			assert.Nil(t, err)
			assert.True(t, p.IsDone())
			assert.True(t, p.IsSuccess())
		})
	}
}

func Test_Empty(t *testing.T) {
	type Case struct {
		name          string
		completeValue int
		shouldError   bool
	}

	cases := []Case{
		{"complete with value", 42, false},
		{"complete with zero", 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Empty[int]()
			assert.False(t, p.IsDone())

			go func() {
				time.Sleep(10 * time.Millisecond)
				if c.shouldError {
					p.CompleteExceptionally(errors.New("test error"))
				} else {
					p.Complete(c.completeValue)
				}
			}()

			val, err := p.Get()
			assert.True(t, p.IsDone())
			assert.Equal(t, c.completeValue, val)
			assert.Equal(t, c.shouldError, err != nil)
		})
	}
}

func Test_Then(t *testing.T) {
	type Case struct {
		name     string
		initial  int
		mapper   func(int) int
		expected int
	}

	cases := []Case{
		{
			name:    "double value",
			initial: 5,
			mapper: func(v int) int {
				return v * 2
			},
			expected: 10,
		},
		{
			name:    "add constant",
			initial: 10,
			mapper: func(v int) int {
				return v + 5
			},
			expected: 15,
		},
		{
			name:    "identity",
			initial: 42,
			mapper: func(v int) int {
				return v
			},
			expected: 42,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Completed(c.initial).Then(c.mapper)
			val, err := p.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Map(t *testing.T) {
	type Case struct {
		name     string
		initial  int
		mapper   func(int) string
		expected string
	}

	cases := []Case{
		{
			name:    "int to string",
			initial: 42,
			mapper: func(v int) string {
				return "value: 42"
			},
			expected: "value: 42",
		},
		{
			name:    "int to empty string",
			initial: 0,
			mapper: func(v int) string {
				return ""
			},
			expected: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Completed(c.initial)
			mapped := promise.Map(p, c.mapper)
			val, err := mapped.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Compose(t *testing.T) {
	type Case struct {
		name     string
		initial  int
		composer func(int) *promise.Promise[string]
		expected string
	}

	cases := []Case{
		{
			name:    "chain two promises",
			initial: 5,
			composer: func(v int) *promise.Promise[string] {
				return promise.Of(func() (string, error) {
					return "result: 5", nil
				})
			},
			expected: "result: 5",
		},
		{
			name:    "chain with completed promise",
			initial: 10,
			composer: func(v int) *promise.Promise[string] {
				return promise.Completed("done")
			},
			expected: "done",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Completed(c.initial)
			composed := promise.Compose(p, c.composer)
			val, err := composed.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_ThenAccept(t *testing.T) {
	type Case struct {
		name    string
		initial int
	}

	cases := []Case{
		{"consume value", 42},
		{"consume zero", 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			consumed := 0
			p := promise.Completed(c.initial)
			voidPromise := p.ThenAccept(func(v int) {
				consumed = v
			})
			_, err := voidPromise.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.initial, consumed)
		})
	}
}

func Test_ThenRun(t *testing.T) {
	type Case struct {
		name string
	}

	cases := []Case{
		{"execute action"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			executed := false
			p := promise.Completed(42)
			voidPromise := p.ThenRun(func() {
				executed = true
			})
			_, err := voidPromise.Get()
			assert.Nil(t, err)
			assert.True(t, executed)
		})
	}
}

func Test_Exceptionally(t *testing.T) {
	type Case struct {
		name        string
		supplier    func() (int, error)
		handler     func(error) int
		expected    int
		shouldError bool
	}

	cases := []Case{
		{
			name: "recover from error",
			supplier: func() (int, error) {
				return 0, errors.New("failed")
			},
			handler: func(err error) int {
				return 99
			},
			expected:    99,
			shouldError: false,
		},
		{
			name: "no error to handle",
			supplier: func() (int, error) {
				return 42, nil
			},
			handler: func(err error) int {
				return 0
			},
			expected:    42,
			shouldError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier).Exceptionally(c.handler)
			val, err := p.Get()
			assert.Equal(t, c.expected, val)
			assert.Equal(t, c.shouldError, err != nil)
		})
	}
}

func Test_Recover(t *testing.T) {
	type Case struct {
		name     string
		supplier func() (string, error)
		handler  func(error) string
		expected string
	}

	cases := []Case{
		{
			name: "recover with fallback",
			supplier: func() (string, error) {
				return "", errors.New("error")
			},
			handler: func(err error) string {
				return "fallback"
			},
			expected: "fallback",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier).Recover(c.handler)
			val, err := p.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Handle(t *testing.T) {
	type Case struct {
		name     string
		supplier func() (int, error)
		handler  func(int, error) string
		expected string
	}

	cases := []Case{
		{
			name: "handle success",
			supplier: func() (int, error) {
				return 42, nil
			},
			handler: func(val int, err error) string {
				if err != nil {
					return "error"
				}
				return "success: 42"
			},
			expected: "success: 42",
		},
		{
			name: "handle error",
			supplier: func() (int, error) {
				return 0, errors.New("failed")
			},
			handler: func(val int, err error) string {
				if err != nil {
					return "error: failed"
				}
				return "success"
			},
			expected: "error: failed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier)
			handled := promise.Handle(p, c.handler)
			val, err := handled.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Join(t *testing.T) {
	type Case struct {
		name        string
		supplier    func() (int, error)
		expected    int
		shouldPanic bool
	}

	cases := []Case{
		{
			name: "join successful promise",
			supplier: func() (int, error) {
				return 42, nil
			},
			expected:    42,
			shouldPanic: false,
		},
		{
			name: "join failed promise",
			supplier: func() (int, error) {
				return 0, errors.New("error")
			},
			expected:    0,
			shouldPanic: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier)
			if c.shouldPanic {
				assert.Panics(t, func() {
					p.Join()
				})
			} else {
				val := p.Join()
				assert.Equal(t, c.expected, val)
			}
		})
	}
}

func Test_Complete(t *testing.T) {
	type Case struct {
		name     string
		value    int
		expected int
		success  bool
	}

	cases := []Case{
		{"complete empty promise", 42, 42, true},
		{"complete with zero", 0, 0, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Empty[int]()
			ok := p.Complete(c.value)
			assert.Equal(t, c.success, ok)
			val, err := p.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_CompleteExceptionally(t *testing.T) {
	type Case struct {
		name    string
		err     error
		success bool
	}

	cases := []Case{
		{"complete with error", errors.New("test error"), true},
		{"complete with nil error", nil, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			promise := promise.Empty[int]()

			ok := promise.CompleteExceptionally(c.err)

			assert.Equal(t, c.success, ok)
		})
	}
}

func Test_Timeout(t *testing.T) {
	type Case struct {
		name        string
		duration    time.Duration
		delay       time.Duration
		shouldError bool
	}

	cases := []Case{
		{
			name:        "timeout before completion",
			duration:    50 * time.Millisecond,
			delay:       200 * time.Millisecond,
			shouldError: true,
		},
		{
			name:        "complete before timeout",
			duration:    200 * time.Millisecond,
			delay:       50 * time.Millisecond,
			shouldError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(func() (int, error) {
				time.Sleep(c.delay)
				return 42, nil
			}).Timeout(c.duration)

			_, err := p.Get()
			assert.Equal(t, c.shouldError, err != nil)
		})
	}
}

func Test_CompleteOnTimeout(t *testing.T) {
	type Case struct {
		name         string
		duration     time.Duration
		delay        time.Duration
		defaultValue int
		expected     int
	}

	cases := []Case{
		{
			name:         "use default on timeout",
			duration:     50 * time.Millisecond,
			delay:        200 * time.Millisecond,
			defaultValue: 99,
			expected:     99,
		},
		{
			name:         "complete before timeout",
			duration:     200 * time.Millisecond,
			delay:        50 * time.Millisecond,
			defaultValue: 99,
			expected:     42,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(func() (int, error) {
				time.Sleep(c.delay)
				return 42, nil
			}).CompleteOnTimeout(c.defaultValue, c.duration)

			val, err := p.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_IsDone(t *testing.T) {
	type Case struct {
		name     string
		promise  *promise.Promise[int]
		expected bool
	}

	cases := []Case{
		{"completed promise", promise.Completed(42), true},
		{"empty promise", promise.Empty[int](), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.promise.IsDone())
		})
	}
}

func Test_IsSuccess(t *testing.T) {
	type Case struct {
		name     string
		supplier func() (int, error)
		expected bool
	}

	cases := []Case{
		{
			name: "successful promise",
			supplier: func() (int, error) {
				return 42, nil
			},
			expected: true,
		},
		{
			name: "failed promise",
			supplier: func() (int, error) {
				return 0, errors.New("error")
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier)
			p.Get()
			assert.Equal(t, c.expected, p.IsSuccess())
		})
	}
}

func Test_IsError(t *testing.T) {
	type Case struct {
		name     string
		supplier func() (int, error)
		expected bool
	}

	cases := []Case{
		{
			name: "error promise",
			supplier: func() (int, error) {
				return 0, errors.New("error")
			},
			expected: true,
		},
		{
			name: "successful promise",
			supplier: func() (int, error) {
				return 42, nil
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier)
			p.Get()
			assert.Equal(t, c.expected, p.IsError())
		})
	}
}

func Test_State(t *testing.T) {
	type Case struct {
		name     string
		promise  *promise.Promise[int]
		expected promise.State
	}

	completedP := promise.Completed(42)
	completedP.Get()

	failedP := promise.Of(func() (int, error) {
		return 0, errors.New("error")
	})
	failedP.Get()

	cases := []Case{
		{"running promise", promise.Empty[int](), promise.StateRunning},
		{"successful promise", completedP, promise.StateSuccess},
		{"failed promise", failedP, promise.StateFailed},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.promise.State())
		})
	}
}

func Test_ChainedOperations(t *testing.T) {
	type Case struct {
		name     string
		initial  int
		expected string
	}

	cases := []Case{
		{"chain multiple operations", 5, "result: 10"},
		{"chain with zero", 0, "result: 0"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Completed(c.initial).Then(func(v int) int {
				return v * 2
			})

			mapped := promise.Map(p, func(v int) string {
				return "result: " + fmt.Sprint(v)
			})

			val, err := mapped.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_ConcurrentCompletion(t *testing.T) {
	type Case struct {
		name     string
		promises int
		expected int
	}

	cases := []Case{
		{"complete multiple promises", 10, 10},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			promises := make([]*promise.Promise[int], c.promises)
			for i := 0; i < c.promises; i++ {
				promises[i] = promise.Of(func() (int, error) {
					time.Sleep(10 * time.Millisecond)
					return 1, nil
				})
			}

			sum := 0
			for _, p := range promises {
				val, err := p.Get()
				assert.Nil(t, err)
				sum += val
			}

			assert.Equal(t, c.expected, sum)
		})
	}
}

func Test_ErrorPropagation(t *testing.T) {
	type Case struct {
		name        string
		initial     func() (int, error)
		shouldError bool
	}

	cases := []Case{
		{
			name: "error propagates through chain",
			initial: func() (int, error) {
				return 0, errors.New("initial error")
			},
			shouldError: true,
		},
		{
			name: "no error in chain",
			initial: func() (int, error) {
				return 42, nil
			},
			shouldError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.initial).
				Then(func(v int) int { return v * 2 })

			mapped := promise.Map(p, func(v int) string {
				return "value"
			})

			_, err := mapped.Get()
			assert.Equal(t, c.shouldError, err != nil)
		})
	}
}
