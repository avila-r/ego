package iterator

import (
	"github.com/avila-r/ego"
)

var errors = ego.ExtendedGoErrorsNamespace.Class("iterator")

var (
	ErrExhausted = errors.New("no more elements in iterator")
)
