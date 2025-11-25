package function

type BiConsumer[T, U any] interface {
	Accept(T, U)
}

type DefaultBiConsumer[T, U any] struct {
	consumer func(T, U)
}

func (d *DefaultBiConsumer[T, U]) Accept(t T, u U) {
	d.consumer(t, u)
}

func NewBiConsumer[T, U any](consumer func(T, U)) BiConsumer[T, U] {
	return &DefaultBiConsumer[T, U]{consumer: consumer}
}

func (d *DefaultBiConsumer[T, U]) AndThenBiConsumer(after BiConsumer[T, U]) BiConsumer[T, U] {
	return NewBiConsumer(func(t T, u U) {
		d.Accept(t, u)
		after.Accept(t, u)
	})
}
