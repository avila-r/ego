package stream

import "github.com/avila-r/ego/function"

type Collector[T, A, R any] interface {
	Supplier() function.Supplier[A]
	Accumulator() function.BiConsumer[A, T]
	Combiner() function.BinaryOperator[A]
	Finisher() function.Function[A, R]
	// TODO: Characteristics() Set<Characteristics>
}

type DefaultCollector[T, A, R any] struct {
	supplier    function.Supplier[A]
	accumulator function.BiConsumer[A, T]
	combiner    function.BinaryOperator[A]
	finisher    function.Function[A, R]
}

func (d *DefaultCollector[T, A, R]) Supplier() function.Supplier[A] {
	return d.supplier
}

func (d *DefaultCollector[T, A, R]) Accumulator() function.BiConsumer[A, T] {
	return d.accumulator
}

func (d *DefaultCollector[T, A, R]) Combiner() function.BinaryOperator[A] {
	return d.combiner
}

func (d *DefaultCollector[T, A, R]) Finisher() function.Function[A, R] {
	return d.finisher
}

func NewCollector[T, A, R any](
	supplier function.Supplier[A],
	accumulator function.BiConsumer[A, T],
	combiner function.BinaryOperator[A],
	finisher function.Function[A, R],
) Collector[T, A, R] {
	return &DefaultCollector[T, A, R]{
		supplier:    supplier,
		accumulator: accumulator,
		combiner:    combiner,
		finisher:    finisher,
	}
}

func (d *DefaultCollector[T, R, R]) Of(
	supplier function.Supplier[R],
	accumulator function.BiConsumer[R, T],
	combiner function.BinaryOperator[R],
	finisher function.Function[R, R],
) Collector[T, R, R] {
	return NewCollector[T, R, R](supplier, accumulator, combiner, finisher)
}

// TODO: Implement Characteristics enum and related methods
