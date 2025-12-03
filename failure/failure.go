package failure

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/avila-r/ego/failure/ctx"
	"github.com/avila-r/ego/failure/property"
	"github.com/avila-r/ego/failure/stacktrace"
	"github.com/avila-r/ego/failure/tags"
	"github.com/avila-r/ego/failure/trait"
)

// Error defines the interface for error handling with extended functionality.
type Error interface {
	Error() string
	Format(state fmt.State, verb rune)
	LogValue() slog.Value
	Class() *ErrorClass
	Message() string
	Cause() error
	Summary() string
	Underlying() []error
	Logs() slog.Value
	Time() time.Time
	Duration() time.Duration
	Domain() string
	Tags() tags.Tags
	Context() ctx.Context
	Trace() string
	Hint() string
	Public() string
	Owner() string
	Span() string
	Trail() string
	Sources() string
	WithOwner(owner string) Error
	WithPublic(public string) Error
	WithHint(hint string) Error
	WithSpan(span string) Error
	WithTrace(trace string) Error
	WithTags(tags tags.Tags) Error
	WithDomain(domain string) Error
	WithDuration(duration time.Duration) Error
	WithDurationSince(t time.Time) Error
	WithTime(time time.Time) Error
	WithCause(err error) Error
	With(key string, value any) Error
	Assert(condition bool, message ...any) Error
	Recover(f func()) error
	Chain() ErrorChain
	Panic()
	Decorate()
	Enhance()
	Decorated() Error
	Enhanced() Error
	Belongs(err error) bool
	Is(err error) bool
	As(target any) bool
	Has(trait trait.Trait) bool
	Extends(c *ErrorClass) bool
	Attribute(key string) property.Result
	Field(key string) property.Result
	Value(key string) property.Result
	Property(key string) property.Result
	Join(errs ...error) Error
	Also(errs ...error) Error
	Unwrap() error
}

func Err(message string, v ...any) error {
	if len(v) > 0 {
		return fmt.Errorf(message, v...)
	}

	return errors.New(message)
}

func Of(message string, v ...any) *Failure {
	return New(message, v...)
}

func New(message string, v ...any) *Failure {
	return Builder(DefaultClass).
		Message(message, v...).
		Build()
}

func Blank() *Failure {
	return Builder(DefaultClass).
		Build()
}

func From(err error) *Failure {
	return Cast(err)
}

func Cast(err error) *Failure {
	if e, ok := err.(*Failure); ok && e != nil {
		return e
	}

	return nil
}

func Decorate(err error, message string, v ...any) *Failure {
	return Builder(transparentWrapper).
		Message(message, v...).
		Cause(err).
		StackTrace().
		Build()
}

func Decorated(err error) *Failure {
	if casted := Cast(err); casted != nil && casted.stacktrace != nil {
		return BuilderFrom(casted).
			StackTrace().
			Build()
	}

	return Builder(stackTraceWrapper).
		Message("").
		Cause(err).
		StackTrace().
		Build()
}

func Enhance(err error, message string, v ...any) *Failure {
	return Builder(transparentWrapper).
		Message(message, v...).
		Cause(err).
		EnhanceStackTrace().
		Build()
}

func Enhanced(err error) *Failure {
	if casted := Cast(err); casted != nil && casted.stacktrace != nil {
		return casted
	}

	return Builder(stackTraceWrapper).
		Message("").
		Cause(err).
		EnhanceStackTrace().
		Build()
}

func Extends(err error, c *ErrorClass) bool {
	casted := func() *Failure {
		raw := err
		for raw != nil {
			typed := Cast(raw)
			if typed != nil {
				return typed
			}
			raw = errors.Unwrap(raw)
		}

		return nil
	}()

	return casted != nil && casted.Extends(c)
}

func Has(err error, trait trait.Trait) bool {
	if casted := Cast(err); casted != nil {
		return casted.Has(trait)
	}

	return false
}

func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	for {
		if reflect.TypeOf(target).Comparable() && err == target {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func As(err error, target any) bool {
	if errors.As(err, &target) {
		return true
	}

	if target == nil || err == nil {
		return false
	}

	val := reflect.ValueOf(target)
	typ := val.Type()

	// target must be a non-nil pointer
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}

	// *target must be interface or implement error
	if typ != reflect.TypeOf(&Failure{}) && reflect.TypeOf(err).AssignableTo(typ.Elem()) {
		return false
	}

	for {
		typeof := reflect.TypeOf(err)
		if typeof != reflect.TypeOf(&Failure{}) && reflect.TypeOf(err).AssignableTo(typ.Elem()) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(interface{ As(any) bool }); ok && x.As(target) {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func Cause(err error) error {
	for {
		previous := Unwrap(err)
		if previous == nil {
			return err
		}
		err = previous
	}
}

func Property(err error, key string) property.Result {
	if err := Cast(err); err != nil {
		return err.Property(key)
	}

	return property.Empty()
}

func Extract[T any](err error, key string) (out T) {
	if err := Cast(err); err != nil {
		err.Property(key).Bind(&out)
	}

	return
}

func Contains(err error, key string) bool {
	if err := Cast(err); err != nil {
		return err.Property(key).Ok
	}

	return false
}

func Inspect(err error) string {
	if err := Cast(err); err != nil {
		return err.Summary()
	}

	return err.Error()
}

func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})

	if !ok {
		return nil
	}

	return u.Unwrap()
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err.Error())
	}
	return v
}

func Try(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = Err("recovered value: %w", v)
			default:
				err = Err("recovered value: %v", r)
			}
		}
	}()
	f()
	return
}

// Panic if error
func Pie(err error) {
	if err != nil {
		panic(err)
	}
}

// Decorated panic if error
func Dpie(err error, message string, v ...any) {
	if err != nil {
		panic(Decorate(err, message, v...))
	}
}

func Deref[T any](p *T, def ...T) (t T) {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	return
}

func Trait(label string) trait.Trait {
	return trait.New(label)
}

func Deep[T comparable](err *Failure, getter func(*Failure) T) T {
	if err.cause == nil {
		return getter(err)
	}

	if casted := Cast(err.cause); casted != nil {
		return func(values ...T) T {
			var zero T
			for _, v := range values {
				if v != zero {
					return v
				}
			}
			return zero
		}(Deep(casted, getter), getter(err))
	}

	return getter(err)
}

func Recurse(err *Failure, do func(*Failure)) {
	do(err)

	if err.cause == nil {
		return
	}

	if casted := Cast(err.cause); casted != nil {
		Recurse(casted, do)
	}
}

func Gather(err *Failure, getter func(*Failure) ctx.Context) ctx.Context {
	if err.cause == nil {
		return getter(err)
	}

	if casted := Cast(err.cause); casted != nil {
		return func(maps ...ctx.Context) ctx.Context {
			count := 0
			for i := range maps {
				count += len(maps[i])
			}

			out := make(ctx.Context, count)
			for i := range maps {
				for k := range maps[i] {
					out[k] = maps[i][k]
				}
			}

			return out
		}(ctx.Context{}, getter(err), Gather(casted, getter))
	}

	return getter(err)
}

func InitializeStackTraceTransformer(subtransformer stacktrace.FilePathTransformer) (stacktrace.FilePathTransformer, error) {
	stacktrace.Transformer.Mu.Lock()
	defer stacktrace.Transformer.Mu.Unlock()

	old := stacktrace.Transformer.Transform.Load().(stacktrace.FilePathTransformer)

	if stacktrace.Transformer.Initialized {
		return old, InitializationFailed.New("stack trace transformer was already set up: %#v", old)
	}

	stacktrace.Transformer.Transform.Store(subtransformer)

	stacktrace.Transformer.Initialized = true
	return nil, nil
}
