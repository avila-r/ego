package result

import "github.com/avila-r/failure"

type Result[T any] struct {
	value *T
	err   error
}

func Of[T any](v T, e error) Result[T] {
	return Result[T]{
		value: &v,
		err:   e,
	}
}

func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: &value,
		err:   nil,
	}
}

func Error[T any](err error) Result[T] {
	return Result[T]{
		value: nil,
		err:   err,
	}
}

func Failure[T any](err error) Result[T] {
	return Error[T](err)
}

func (r *Result[T]) Value() *T {
	return r.value
}

func (r *Result[T]) Error() error {
	if r.err == nil && r.IsEmpty() {
		return ErrEmptyResult
	}
	return r.err
}

func (r *Result[T]) Success(t T) Result[T] {
	r.value, r.err = &t, nil
	return *r
}

func (r *Result[T]) Ok(t T) Result[T] {
	r.value, r.err = &t, nil
	return *r
}

func (r *Result[T]) Failure(err error) Result[T] {
	r.err, r.value = err, nil
	return *r
}

func (r *Result[T]) IsEmpty() bool {
	return r.value == nil
}

func (r *Result[T]) IsSuccess() bool {
	return r.Error() == nil && r.value != nil
}

func (r *Result[T]) IsError() bool {
	return r.Error() != nil || r.value == nil
}

func (r *Result[T]) Unwrap() T {
	return *r.value
}

func (o *Result[T]) Take() (*T, *failure.Error) {
	if o.IsEmpty() {
		println("taking from empty result")
		return nil, ErrNoneValueTaken
	}
	return o.value, nil
}

func (r *Result[T]) Join() T {
	if r.IsError() {
		panic(r.Error())
	}
	return *r.value
}

func (r *Result[T]) Expect(message ...string) T {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = r.Error().Error()
	}

	if r.IsError() {
		panic(msg)
	}

	return *r.value
}

// Map applies a function to the value if Ok, otherwise returns the error
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsError() {
		return r
	}
	return Ok(f(*r.value))
}

// MapOr applies a function to the value if Ok, otherwise returns the default value
func (r Result[T]) MapOr(defaultValue T, f func(T) T) T {
	if r.IsError() {
		return defaultValue
	}
	return f(*r.value)
}

// MapOrElse applies a function to the value if Ok, otherwise calls the fallback function
func (r *Result[T]) MapOrElse(fallback func(error) T, f func(T) T) T {
	if r.IsError() {
		return fallback(r.err)
	}
	return f(*r.value)
}

// FlatMap applies a function that returns a Result and flattens it
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T] {
	if r.IsError() {
		return r
	}
	return f(*r.value)
}

// AndThen chains operations that return Results (alias for FlatMap)
func (r Result[T]) AndThen(f func(T) Result[T]) Result[T] {
	return r.FlatMap(f)
}

// Bind is another alias for FlatMap, common in functional programming
func (r Result[T]) Bind(f func(T) Result[T]) Result[T] {
	return r.FlatMap(f)
}

// OnSuccess executes a side-effect function if the result is Ok
func (r Result[T]) OnSuccess(f func(T)) Result[T] {
	if r.IsSuccess() {
		f(*r.value)
	}
	return r
}

// OnFailure executes a side-effect function if the result is Error
func (r Result[T]) OnFailure(f func(error)) Result[T] {
	if r.IsError() {
		f(r.err)
	}
	return r
}

// Inspect allows inspecting the value without consuming it (side-effect)
func (r Result[T]) Inspect(f func(T)) Result[T] {
	return r.OnSuccess(f)
}

// InspectErr allows inspecting the error without consuming it
func (r Result[T]) InspectErr(f func(error)) Result[T] {
	return r.OnFailure(f)
}

// Or returns this result if Ok, otherwise returns the other result
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.IsSuccess() {
		return r
	}
	return other
}

// OrElse returns this result if Ok, otherwise calls the function
func (r Result[T]) OrElse(f func(error) Result[T]) Result[T] {
	if r.IsSuccess() {
		return r
	}
	return f(r.err)
}

// And returns other if this is Ok, otherwise returns this error
func (r *Result[T]) And(other Result[T]) Result[T] {
	if r.IsError() {
		return *r
	}
	return other
}

// UnwrapOr returns the value if Ok, otherwise returns the default value
func (r *Result[T]) UnwrapOr(defaultValue T) T {
	if r.IsError() {
		return defaultValue
	}
	return *r.value
}

// UnwrapOrElse returns the value if Ok, otherwise calls the function
func (r Result[T]) UnwrapOrElse(f func(error) T) T {
	if r.IsError() {
		return f(r.err)
	}
	return *r.value
}

// UnwrapOrDefault returns the value if Ok, otherwise returns the zero value
func (r Result[T]) UnwrapOrDefault() T {
	var zero T
	return r.UnwrapOr(zero)
}

// Contains checks if the result contains the given value
func (r Result[T]) Contains(value T, equals func(T, T) bool) bool {
	if r.IsError() {
		return false
	}
	return equals(*r.value, value)
}

// ContainsErr checks if the result contains an error matching the predicate
func (r *Result[T]) ContainsErr(predicate func(error) bool) bool {
	if r.IsSuccess() {
		return false
	}
	return predicate(r.err)
}

// Match performs pattern matching on the result
func (r Result[T]) Match(onSuccess func(T), onFailure func(error)) {
	if r.IsSuccess() {
		onSuccess(*r.value)
	} else {
		onFailure(r.err)
	}
}

// MatchReturn performs pattern matching and returns a value
func MatchReturn[T any, U any](r Result[T], onSuccess func(T) U, onFailure func(error) U) U {
	if r.IsSuccess() {
		return onSuccess(*r.value)
	}
	return onFailure(r.err)
}

// Filter keeps the value only if the predicate is true, otherwise returns error
func (r Result[T]) Filter(predicate func(T) bool, err error) Result[T] {
	if r.IsError() {
		return r
	}
	if predicate(*r.value) {
		return r
	}
	return Error[T](err)
}

// Recover attempts to recover from an error using the provided function
func (r Result[T]) Recover(f func(error) T) Result[T] {
	if r.IsSuccess() {
		return r
	}
	return Ok(f(r.err))
}

// RecoverWith attempts to recover from an error with a Result
func (r Result[T]) RecoverWith(f func(error) Result[T]) Result[T] {
	if r.IsSuccess() {
		return r
	}
	return f(r.err)
}

// Flatten flattens a Result[Result[T]] into Result[T]
func Flatten[T any](r Result[Result[T]]) Result[T] {
	if r.IsError() {
		return Error[T](r.err)
	}
	return *r.value
}

// Transpose converts Result[*T] to *Result[T] style handling
func (r *Result[T]) Transpose() (*T, error) {
	if r.IsError() {
		return nil, r.err
	}
	return r.value, nil
}

// Try executes a function and captures panics as errors
func Try[T any](f func() T) Result[T] {
	var result Result[T]
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				result = Error[T](err)
			} else {
				result = Error[T](failure.New("panic occurred"))
			}
		}
	}()
	result = Ok(f())
	return result
}

// TryWith executes a fallible function and wraps it in a Result
func TryWith[T any](f func() (T, error)) Result[T] {
	v, err := f()
	return Of(v, err)
}

// Tap executes a function for side effects and returns the original result
func (r Result[T]) Tap(f func(Result[T])) Result[T] {
	f(r)
	return r
}

// ToPointer converts Result[T] to traditional Go (*T, error) pattern
func (r *Result[T]) ToPointer() (*T, error) {
	return r.Transpose()
}

// FromPointer creates a Result from traditional Go (*T, error) pattern
func FromPointer[T any](value *T, err error) Result[T] {
	if err != nil {
		return Error[T](err)
	}
	if value == nil {
		return Error[T](ErrEmptyResult)
	}
	return Ok(*value)
}

// Chain allows chaining multiple Result operations
func Chain[T any](initial Result[T], operations ...func(Result[T]) Result[T]) Result[T] {
	result := initial
	for _, op := range operations {
		result = op(result)
		if result.IsError() {
			return result
		}
	}
	return result
}

// Combine combines two Results into a Result of a tuple-like struct
func Combine[T any, U any](r1 Result[T], r2 Result[U]) Result[struct {
	First  T
	Second U
}] {
	if r1.IsError() {
		return Error[struct {
			First  T
			Second U
		}](r1.err)
	}
	if r2.IsError() {
		return Error[struct {
			First  T
			Second U
		}](r2.err)
	}
	return Ok(struct {
		First  T
		Second U
	}{
		First:  *r1.value,
		Second: *r2.value,
	})
}

// MapErr transforms the error if present
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.IsSuccess() {
		return r
	}
	return Error[T](f(r.err))
}

// OkOr returns the value if Ok, otherwise returns error as Result
func (r Result[T]) OkOr(err error) Result[T] {
	if r.IsSuccess() {
		return r
	}
	return Error[T](err)
}

// ExpectErr panics with a message if the result is Ok (useful for testing)
func (r *Result[T]) ExpectErr(message string) error {
	if r.IsSuccess() {
		panic(message)
	}
	return r.err
}

// UnwrapErr returns the error, panics if Ok
func (r *Result[T]) UnwrapErr() error {
	if r.IsSuccess() {
		panic("called UnwrapErr on an Ok value")
	}
	return r.err
}

// Iter returns a slice with the value if Ok, empty slice if Error
func (r *Result[T]) Iter() []T {
	if r.IsError() {
		return []T{}
	}
	return []T{*r.value}
}

// Collect gathers multiple Results into a single Result of a slice
func Collect[T any](results []Result[T]) Result[[]T] {
	var values []T
	for _, r := range results {
		if r.IsError() {
			return Error[[]T](r.err)
		}
		values = append(values, *r.value)
	}
	return Ok(values)
}

// Partition separates a slice of Results into successes and failures
func Partition[T any](results []Result[T]) (successes []T, failures []error) {
	for _, r := range results {
		if r.IsSuccess() {
			successes = append(successes, *r.value)
		} else {
			failures = append(failures, r.err)
		}
	}
	return successes, failures
}
