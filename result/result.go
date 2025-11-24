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

func Err[T any](err error) Result[T] {
	return Error[T](err)
}

func (r Result[T]) Value() *T {
	return r.value
}

func (r Result[T]) Error() error {
	if r.err == nil && r.IsEmpty() {
		return ErrEmptyResult
	}
	return r.err
}

func (r Result[T]) IsEmpty() bool {
	return r.value == nil
}

func (r Result[T]) IsSuccess() bool {
	return r.Error() == nil && r.value != nil
}

func (r Result[T]) IsError() bool {
	return r.Error() != nil || r.value == nil
}

func (r Result[T]) Unwrap() T {
	return *r.value
}

func (o Result[T]) Take() (*T, *failure.Error) {
	if o.IsEmpty() {
		return nil, ErrNoneValueTaken
	}
	return o.value, nil
}

func (r Result[T]) Join() T {
	if r.IsError() {
		panic(r.Error())
	}
	return *r.value
}

func (r Result[T]) Expect(message ...string) T {
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
