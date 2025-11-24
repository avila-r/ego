package optional

import (
	"github.com/avila-r/ego"
)

var errors = ego.ExtendedGoErrorsNamespace.Class("optional")

var (
	ErrNoneValueTaken = errors.New("none value taken")
	ErrNoPresentValue = errors.New("value is not present")
)
