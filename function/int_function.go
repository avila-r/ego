package function

type IntFunction interface {
	Apply(int) int
}

type DefaultIntFunction struct {
	function func(int) int
}

func (d *DefaultIntFunction) Apply(t int) int {
	return d.function(t)
}

func NewIntFunction(function func(int) int) IntFunction {
	return &DefaultIntFunction{function: function}
}
