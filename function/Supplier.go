package function

type Supplier[T any] interface {
	Get() T
}

type DefaultSupplier[T any] struct {
	supplier func() T
}

func (d *DefaultSupplier[T]) Get() T {
	return d.supplier()
}

func NewSupplier[T any](supplier func() T) Supplier[T] {
	return &DefaultSupplier[T]{supplier: supplier}
}
