package function

type BinaryOperator[T any] interface {
	BiFunction[T, T, T]
	Apply(T, T) T
}

type DefaultBinaryOperator[T any] struct {
	operator func(T, T) T
}

func (d *DefaultBinaryOperator[T]) Apply(t1 T, t2 T) T {
	return d.operator(t1, t2)
}

func NewBinaryOperator[T any](operator func(T, T) T) BinaryOperator[T] {
	return &DefaultBinaryOperator[T]{operator: operator}
}

func (d *DefaultBinaryOperator[T]) MinBy(comparator Comparator[T]) BinaryOperator[T] {
	return NewBinaryOperator(func(t1 T, t2 T) T {
		if comparator.Compare(t1, t2) <= 0 {
			return t1
		}
		return t2
	})
}

func (d *DefaultBinaryOperator[T]) MaxBy(comparator Comparator[T]) BinaryOperator[T] {
	return NewBinaryOperator(func(t1 T, t2 T) T {
		if comparator.Compare(t1, t2) >= 0 {
			return t1
		}
		return t2
	})
}
