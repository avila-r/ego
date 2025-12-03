package optional

import "github.com/avila-r/ego/failure"

type Optional[T any] struct {
	value *T
}

func Of[T any](value T) Optional[T] {
	return Optional[T]{value: &value}
}

func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

func (o *Optional[T]) IsPresent() bool {
	return o.value != nil
}

func (o *Optional[T]) IsEmpty() bool {
	return o.value == nil
}

func (o *Optional[T]) Join() T {
	if !o.IsPresent() {
		panic(ErrNoPresentValue)
	}
	return *o.value
}

func (o *Optional[T]) Get() (t T, ok bool) {
	if o.IsEmpty() {
		return
	}
	return *o.value, true
}

func (o *Optional[T]) Clear() {
	o.value = nil
}

func (o *Optional[T]) GetOrDefault(fallback T) T {
	if o.IsPresent() {
		return *o.value
	}
	return fallback
}

func (o *Optional[T]) Set(value T) {
	o.value = &value
}

func (o *Optional[T]) Take() (*T, failure.Error) {
	if o.IsEmpty() {
		return nil, ErrNoneValueTaken
	}
	return o.value, nil
}
