package function

type Consumer[T any] interface {
	Accept(T)
}

type DefaultConsumer[T any] struct {
	consumer func(T)
}

func (d *DefaultConsumer[T]) Accept(t T) {
	d.consumer(t)
}

func NewConsumer[T any](consumer func(T)) Consumer[T] {
	return &DefaultConsumer[T]{consumer: consumer}
}

func (d *DefaultConsumer[T]) AndThenConsumer(after Consumer[T]) Consumer[T] {
	return NewConsumer(func(t T) {
		d.Accept(t)
		after.Accept(t)
	})
}
