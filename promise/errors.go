package promise

import "github.com/avila-r/ego"

var errors = ego.ExtendedGoErrorsNamespace.Class("promise")

var (
	ErrTimeout   = errors.New("promise timeout")
	ErrCancelled = errors.New("promise cancelled")
)
