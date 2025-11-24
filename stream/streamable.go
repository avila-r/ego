package stream

type Streamable[T comparable] interface {
	Stream() Stream[T]
}

type Collectable[T any] interface {
	Elements() []T
}
