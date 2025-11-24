package box

// Of creates a new Box containing the given value
func Of[T any](value T) *Box[T] {
	return &Box[T]{
		value:   &value,
		present: true,
	}
}

// Empty creates a new empty Box
func Empty[T any]() *Box[T] {
	return &Box[T]{
		value:   nil,
		present: false,
	}
}

// Integer helper functions
func fromInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// Increment returns a function that increments an integer by 1
func Increment() func(*int) int {
	return func(v *int) int {
		return fromInt(v) + 1
	}
}

// Decrement returns a function that decrements an integer by 1
func Decrement() func(*int) int {
	return func(v *int) int {
		return fromInt(v) - 1
	}
}

// Square returns a function that squares an integer
func Square() func(*int) int {
	return func(v *int) int {
		x := fromInt(v)
		return x * x
	}
}

// Cube returns a function that cubes an integer
func Cube() func(*int) int {
	return func(v *int) int {
		x := fromInt(v)
		return x * x * x
	}
}

// Twice returns a function that doubles an integer
func Twice() func(*int) int {
	return func(v *int) int {
		return fromInt(v) * 2
	}
}

// Halve returns a function that halves an integer
func Halve() func(*int) int {
	return func(v *int) int {
		return fromInt(v) / 2
	}
}

// Negate returns a function that negates an integer
func Negate() func(*int) int {
	return func(v *int) int {
		return -fromInt(v)
	}
}

// Abs returns a function that returns the absolute value of an integer
func Abs() func(*int) int {
	return func(v *int) int {
		x := fromInt(v)
		if x < 0 {
			return -x
		}
		return x
	}
}

// Identity returns a function that returns the integer as-is
func Identity() func(*int) int {
	return fromInt
}

// Modulo returns a function that computes v % m
func Modulo(m int) func(*int) int {
	return func(v *int) int {
		return fromInt(v) % m
	}
}

// Clamp returns a function that clamps an integer between min and max
func Clamp(min, max int) func(*int) int {
	return func(v *int) int {
		x := fromInt(v)
		if x < min {
			return min
		}
		if x > max {
			return max
		}
		return x
	}
}
