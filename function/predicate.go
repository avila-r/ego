package function

type Predicate[T any] interface {
	Test(T) bool
}

type DefaultPredicate[T any] struct {
	predicate func(T) bool
}

func (d *DefaultPredicate[T]) Test(t T) bool {
	return d.predicate(t)
}

func NewPredicate[T any](predicate func(T) bool) Predicate[T] {
	return &DefaultPredicate[T]{predicate: predicate}
}

func (d *DefaultPredicate[T]) And(other Predicate[T]) Predicate[T] {
	return NewPredicate(func(t T) bool {
		return d.Test(t) && other.Test(t)
	})
}

func (d *DefaultPredicate[T]) Or(other Predicate[T]) Predicate[T] {
	return NewPredicate(func(t T) bool {
		return d.Test(t) || other.Test(t)
	})
}

func (d *DefaultPredicate[T]) Negate() Predicate[T] {
	return NewPredicate(func(t T) bool {
		return !d.Test(t)
	})
}

func IsEqual[T comparable](target T) Predicate[T] {
	return NewPredicate(func(t T) bool {
		return t == target
	})
}

func Not[T any](predicate DefaultPredicate[T]) Predicate[T] {
	return predicate.Negate()
}
