package function

type ToFloatFunction[T any] interface {
	ApplyAsFloat(T) float64
}

type DefaultToFloatFunction[T any] struct {
	function func(T) float64
}

func (d *DefaultToFloatFunction[T]) ApplyAsFloat(t T) float64 {
	return d.function(t)
}

func NewToFloatFunction[T any](function func(T) float64) ToFloatFunction[T] {
	return &DefaultToFloatFunction[T]{function: function}
}
