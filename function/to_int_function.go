package function

type ToIntFunction[T any] interface {
	ApplyAsInt(T) int
}

type DefaultToIntFunction[T any] struct {
	function func(T) int
}

func (d *DefaultToIntFunction[T]) ApplyAsInt(t T) int {
	return d.function(t)
}

func NewToIntFunction[T any](function func(T) int) ToIntFunction[T] {
	return &DefaultToIntFunction[T]{function: function}
}
