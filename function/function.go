package function

type Function[T, R any] interface {
	Apply(T) R
}

type DefaultFunction[T, R any] struct {
	function func(T) R
}

func (d *DefaultFunction[T, R]) Apply(t T) R {
	return d.function(t)
}

func NewFunction[T, R any](function func(T) R) Function[T, R] {
	return &DefaultFunction[T, R]{function: function}
}

func ComposeFunction[V, T, R any](this Function[T, R], before Function[V, T]) Function[V, R] {
	return NewFunction(func(v V) R {
		t := before.Apply(v)
		return this.Apply(t)
	})
}

func AndThenFunction[T, R, V any](this Function[T, R], after Function[R, V]) Function[T, V] {
	return NewFunction(func(t T) V {
		r := this.Apply(t)
		return after.Apply(r)
	})
}

func IdentityFunction[T any]() Function[T, T] {
	return NewFunction(func(t T) T {
		return t
	})
}
