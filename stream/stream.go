package stream

type Stream[T comparable] struct {
	elements []T
}

func Of[T comparable](elements ...T) Stream[T] {
	return Stream[T]{elements: elements}
}

func From[T comparable](collectable Collectable[T]) Stream[T] {
	return Stream[T]{elements: collectable.Elements()}
}
