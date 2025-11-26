package box

var (
	// Increment increments an integer by 1
	Increment = func(v int) int {
		return v + 1
	}

	// Decrement decrements an integer by 1
	Decrement = func(v int) int {
		return v - 1
	}

	// Square squares an integer
	Square = func(v int) int {
		return v * v
	}

	// Cube cubes an integer
	Cube = func(v int) int {
		return v * v * v
	}

	// Twice doubles an integer
	Twice = func(v int) int {
		return v * 2
	}

	// Halve halves an integer
	Halve = func(v int) int {
		return v / 2
	}

	// Negate negates an integer
	Negate = func(v int) int {
		return -v
	}

	// Abs returns the absolute value of an integer
	Abs = func(v int) int {
		if v < 0 {
			return -v
		}
		return v
	}

	// Identity returns the same integer
	Identity = func(v int) int {
		return v
	}

	// Modulo returns v % m
	Modulo = func(m int) func(int) int {
		return func(v int) int {
			return v % m
		}
	}

	// Clamp clamps v between min and max
	Clamp = func(min, max int) func(int) int {
		return func(v int) int {
			if v < min {
				return min
			}
			if v > max {
				return max
			}
			return v
		}
	}
)
