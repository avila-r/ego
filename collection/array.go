package collection

import (
	"github.com/avila-r/ego/constraint"
	"github.com/avila-r/ego/stream"
)

type Array[T constraint.Comparable] interface {
	stream.Streamable[T]
}
