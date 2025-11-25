package function

type UnaryOperator[T any] interface {
	Function[T, T]
}

type DefaultUnaryOperator[T any] struct {
	function func(T) T
}

func (d *DefaultUnaryOperator[T]) Apply(t T) T {
	return d.function(t)
}

func NewUnaryOperator[T any](function func(T) T) UnaryOperator[T] {
	return &DefaultUnaryOperator[T]{function: function}
}

func IdentityUnaryOperator[T any]() UnaryOperator[T] {
	return NewUnaryOperator(func(t T) T {
		return t
	})
}
