package collection

import (
	"github.com/avila-r/ego/iterator"
	"github.com/avila-r/ego/stream"
)

type List[T comparable] interface {
	Add(elements ...T)
	Get(index int) (T, bool)
	Set(index int, element T) bool
	Remove(index int) bool
	Size() int
	IsEmpty() bool
	Clear()
	Contains(element T) bool

	stream.Collectable[T]
	stream.Streamable[T]
	iterator.Iterable[T]
}
