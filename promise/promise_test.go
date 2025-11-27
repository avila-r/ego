package promise_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/avila-r/ego/pointer"
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

func Test_Supply(t *testing.T) {
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
			p := promise.Supply(c.supplier)
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

func Test_Done(t *testing.T) {
	type Case struct {
		name     string
		value    string
		expected string
	}

	cases := []Case{
		{"already done", "hello", "hello"},
		{"empty string", "", ""},
		{"number as string", "123", "123"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Done(c.value)
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

func Test_ThenWithConcurrency(t *testing.T) {
	type Case struct {
		name        string
		initial     int
		concurrency promise.Concurrency
		mapper      func(int) int
		expected    int
	}

	cases := []Case{
		{
			name:    "double value",
			initial: 5,
			mapper: func(v int) int {
				return v * 2
			},
			concurrency: promise.Async,
			expected:    10,
		},
		{
			name:    "add constant",
			initial: 10,
			mapper: func(v int) int {
				return v + 5
			},
			concurrency: promise.Sync,
			expected:    15,
		},
		{
			name:    "identity",
			initial: 42,
			mapper: func(v int) int {
				return v
			},
			concurrency: promise.Async,
			expected:    42,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Completed(c.initial).ThenWithConcurrency(c.concurrency, c.mapper)
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

func Test_MapWithConcurrency(t *testing.T) {
	type Case struct {
		name        string
		initial     *int
		mapper      func(int) string
		expected    string
		concurrency promise.Concurrency
	}

	cases := []Case{
		{
			name:    "int to string async",
			initial: pointer.Of(42),
			mapper: func(v int) string {
				return "value: 42"
			},
			expected:    "value: 42",
			concurrency: promise.Async,
		},
		{
			name:    "int to string sync",
			initial: pointer.Of(42),
			mapper: func(v int) string {
				return "value: 42"
			},
			expected:    "value: 42",
			concurrency: promise.Sync,
		},
		{
			name:    "int to empty string sync",
			initial: pointer.Of(0),
			mapper: func(v int) string {
				return ""
			},
			expected:    "",
			concurrency: promise.Sync,
		},
		{
			name:    "failure sync",
			initial: nil,
			mapper: func(v int) string {
				return ""
			},
			expected:    "",
			concurrency: promise.Sync,
		},
		{
			name:    "failure async",
			initial: nil,
			mapper: func(v int) string {
				return ""
			},
			expected:    "",
			concurrency: promise.Async,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(func() (int, error) {
				if c.initial == nil {
					return 0, errors.New("not completed")
				}
				return *c.initial, nil
			})

			mapped := promise.MapWithConcurrency(p, c.concurrency, c.mapper)
			val, err := mapped.Get()
			if c.initial == nil {
				assert.NotNil(t, err)
				return
			}
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

func Test_ComposeWithConcurrency(t *testing.T) {
	type Case struct {
		name        string
		initial     *int
		composer    func(int) *promise.Promise[string]
		expected    string
		concurrency promise.Concurrency
	}

	cases := []Case{
		{
			name:    "chain two promises (sync)",
			initial: pointer.Of(5),
			composer: func(v int) *promise.Promise[string] {
				return promise.Of(func() (string, error) {
					return "result: 5", nil
				})
			},
			expected:    "result: 5",
			concurrency: promise.Sync,
		},
		{
			name:    "chain with completed promise (async)",
			initial: pointer.Of(10),
			composer: func(v int) *promise.Promise[string] {
				return promise.Completed("done")
			},
			expected:    "done",
			concurrency: promise.Async,
		},
		{
			name:    "error (async)",
			initial: nil,
			composer: func(v int) *promise.Promise[string] {
				return promise.Of(func() (string, error) {
					return "", errors.New("not completed")
				})
			},
			expected:    "",
			concurrency: promise.Async,
		},
		{
			name:    "error (sync)",
			initial: nil,
			composer: func(v int) *promise.Promise[string] {
				return promise.Of(func() (string, error) {
					return "", errors.New("not completed")
				})
			},
			expected:    "",
			concurrency: promise.Sync,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(func() (int, error) {
				if c.initial == nil {
					return 0, errors.New("not completed")
				}
				return *c.initial, nil
			})

			composed := promise.ComposeWithConcurrency(p, c.concurrency, c.composer)

			val, err := composed.Get()
			if c.initial == nil {
				println("initial is nil")
				assert.Empty(t, val)
				assert.NotNil(t, err)
				return
			}

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

func Test_ThenAcceptWithConcurrency(t *testing.T) {
	type Case struct {
		name        string
		initial     *int
		concurrency promise.Concurrency
	}

	cases := []Case{
		{"consume value", pointer.Of(42), promise.Async},
		{"consume zero", pointer.Of(0), promise.Sync},
		{"uncompleted async", nil, promise.Async},
		{"uncompleted sync", nil, promise.Sync},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			consumed := 0

			p := promise.Of(func() (int, error) {
				if c.initial == nil {
					return 0, errors.New("not completed")
				}

				return *c.initial, nil
			})

			void := p.ThenAcceptWithConcurrency(c.concurrency, func(v int) {
				consumed = v
			})

			_, err := void.Get()

			if c.initial == nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, c.initial, pointer.Of(consumed))
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

func Test_ThenRunWithConcurrency(t *testing.T) {
	type Case struct {
		name        string
		concurrency promise.Concurrency
	}

	cases := []Case{
		{"execute action sync", promise.Sync},
		{"execute action async", promise.Async},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			executed := false

			p := promise.Completed(42)

			void := p.ThenRunWithConcurrency(c.concurrency, func() {
				executed = true
			})

			_, err := void.Get()

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

func Test_ExceptionallyWithConcurrency(t *testing.T) {
	type Case struct {
		name        string
		supplier    func() (int, error)
		handler     func(error) int
		expected    int
		shouldError bool
		concurrency promise.Concurrency
	}

	cases := []Case{
		{
			name: "recover from error (async)",
			supplier: func() (int, error) {
				return 0, errors.New("failed")
			},
			handler: func(err error) int {
				return 99
			},
			expected:    99,
			shouldError: false,
			concurrency: promise.Async,
		},
		{
			name: "no error to handle (async)",
			supplier: func() (int, error) {
				return 42, nil
			},
			handler: func(err error) int {
				return 0
			},
			expected:    42,
			shouldError: false,
			concurrency: promise.Async,
		},
		{
			name: "recover from error (sync)",
			supplier: func() (int, error) {
				return 0, errors.New("failed")
			},
			handler: func(err error) int {
				return 99
			},
			expected:    99,
			shouldError: false,
			concurrency: promise.Sync,
		},
		{
			name: "no error to handle (sync)",
			supplier: func() (int, error) {
				return 42, nil
			},
			handler: func(err error) int {
				return 0
			},
			expected:    42,
			shouldError: false,
			concurrency: promise.Sync,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier).ExceptionallyWithConcurrency(c.concurrency, c.handler)
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

func Test_HandleWithConcurrency(t *testing.T) {
	type Case struct {
		name        string
		supplier    func() (int, error)
		handler     func(int, error) string
		expected    string
		concurrency promise.Concurrency
	}

	cases := []Case{
		{
			name: "handle success (async)",
			supplier: func() (int, error) {
				return 42, nil
			},
			handler: func(val int, err error) string {
				if err != nil {
					return "error"
				}
				return "success: 42"
			},
			expected:    "success: 42",
			concurrency: promise.Async,
		},
		{
			name: "handle error (async)",
			supplier: func() (int, error) {
				return 0, errors.New("failed")
			},
			handler: func(val int, err error) string {
				if err != nil {
					return "error: failed"
				}
				return "success"
			},
			expected:    "error: failed",
			concurrency: promise.Async,
		},
		{
			name: "handle success (sync)",
			supplier: func() (int, error) {
				return 42, nil
			},
			handler: func(val int, err error) string {
				if err != nil {
					return "error"
				}
				return "success: 42"
			},
			expected:    "success: 42",
			concurrency: promise.Sync,
		},
		{
			name: "handle error (sync)",
			supplier: func() (int, error) {
				return 0, errors.New("failed")
			},
			handler: func(val int, err error) string {
				if err != nil {
					return "error: failed"
				}
				return "success"
			},
			expected:    "error: failed",
			concurrency: promise.Sync,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := promise.Of(c.supplier)
			handled := promise.HandleWithConcurrency(p, c.concurrency, c.handler)
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
		promise  *promise.Promise[int]
		name     string
		value    int
		expected int
		success  bool
	}

	cases := []Case{
		{promise.Empty[int](), "complete empty promise", 42, 42, true},
		{promise.Empty[int](), "complete with zero", 0, 0, true},
		{promise.Completed(2), "already completed", 32, 2, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ok := c.promise.Complete(c.value)
			assert.Equal(t, c.success, ok)
			val, err := c.promise.Get()
			assert.Nil(t, err)
			assert.Equal(t, c.expected, val)
		})
	}
}

func Test_Cancel(t *testing.T) {
	type Case struct {
		name        string
		promise     *promise.Promise[int]
		success     bool
		expectState promise.State
	}

	cases := []Case{
		{
			name:        "cancel running promise",
			promise:     promise.Empty[int](),
			success:     true,
			expectState: promise.StateCancelled,
		},
		{
			name: "cancel already cancelled",
			promise: func() *promise.Promise[int] {
				p := promise.Empty[int]()
				p.Cancel()
				return p
			}(),
			success:     false,
			expectState: promise.StateCancelled,
		},
		{
			name: "cancel completed promise has no effect",
			promise: func() *promise.Promise[int] {
				p := promise.Empty[int]()
				p.Complete(10)
				return p
			}(),
			success:     false,
			expectState: promise.StateCompleted,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// capture old state before cancelling
			oldState := c.promise.State()

			c.promise.Cancel()

			now := c.promise.State()

			if c.success {
				assert.Equal(t, c.expectState, now)
			} else {
				// if cancel was not supposed to have effect, state stays as it was
				assert.Equal(t, oldState, now)
			}

			// check Get() behavior
			_, err := c.promise.Get()

			if now == promise.StateCancelled {
				assert.Error(t, err, "Get() should fail for cancelled promise")
				assert.True(t, c.promise.IsCancelled())
			} else {
				assert.NoError(t, err, "Get() should succeed when not cancelled")
			}
		})
	}
}

func Test_CompleteExceptionally(t *testing.T) {
	type Case struct {
		promise *promise.Promise[int]
		name    string
		err     error
		success bool
	}

	cases := []Case{
		{promise.Empty[int](), "complete with error", errors.New("test error"), true},
		{promise.Empty[int](), "complete with nil error", nil, true},
		{promise.Completed(42), "already completed", errors.New("test error"), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ok := c.promise.CompleteExceptionally(c.err)

			assert.Equal(t, c.success, ok)
			if c.success {
				assert.True(t, c.promise.IsCompletedExceptionally())
			}
		})
	}
}

func Test_Context(t *testing.T) {
	t.Run("Context is not nil", func(t *testing.T) {
		p := promise.Empty[int]()
		ctx := p.Context()
		assert.NotNil(t, ctx, "Context should not be nil")
	})

	t.Run("Context reflects cancellation", func(t *testing.T) {
		p, cancel := promise.WithCancel[int]()

		defer cancel()

		ctx := p.Context()

		assert.NotNil(t, ctx)

		cancel()

		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
			t.Error("Expected context to be cancelled")
		}
	})
}

func Test_State_String(t *testing.T) {
	cases := []struct {
		state    promise.State
		expected string
	}{
		{promise.StateRunning, "Running"},
		{promise.StateCompleted, "Success"},
		{promise.StateFailed, "Failed"},
		{promise.StateCancelled, "Cancelled"},
	}

	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			assert.Equal(t, c.expected, c.state.String())
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

func Test_Threshold(t *testing.T) {
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
			}).Threshold(c.duration)

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

func Test_IsCompleted(t *testing.T) {
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
			assert.Equal(t, c.expected, c.promise.IsCompleted())
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
		{"successful promise", completedP, promise.StateCompleted},
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
