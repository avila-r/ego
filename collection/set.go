package collection

import (
	"github.com/avila-r/ego/stream"
)

type Set[T comparable] interface {
	stream.Streamable[T]
}
