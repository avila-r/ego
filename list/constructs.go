package list

import (
	"github.com/avila-r/ego/collection"
)

func New[T comparable]() collection.List[T] {
	return EmptyArrayList[T]()
}

func Empty[T comparable]() collection.List[T] {
	return New[T]()
}

func Of[T comparable](elements ...T) collection.List[T] {
	list := New[T]()
	list.Add(elements...)
	return list
}
