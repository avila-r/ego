package function

type BiFunction[T, U, R any] interface {
	Apply(T, U) R
}

type DefaultBiFunction[T, U, R any] struct {
	function func(T, U) R
}

func (d *DefaultBiFunction[T, U, R]) Apply(t T, u U) R {
	return d.function(t, u)
}

func NewBiFunction[T, U, R any](function func(T, U) R) BiFunction[T, U, R] {
	return &DefaultBiFunction[T, U, R]{function: function}
}

func AndThenBiFunction[T, U, R, V any](this BiFunction[T, U, R], after BiFunction[R, U, V]) BiFunction[T, U, V] {
	return NewBiFunction(func(t T, u U) V {
		return after.Apply(this.Apply(t, u), u)
	})
}
