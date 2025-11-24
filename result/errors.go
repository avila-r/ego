package result

import (
	"github.com/avila-r/ego"
)

var errors = ego.ExtendedGoErrorsNamespace.Class("result")

var (
	ErrEmptyResult    = errors.New("empty result - unexpected nil value and error")
	ErrNoneValueTaken = errors.New("none value taken")
	ErrNoPresentValue = errors.New("value is not present")
)
