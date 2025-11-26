package box

// Of creates a new Box containing the given value
func Of[T any](value T) Box[T] {
	return Box[T]{
		value:   &value,
		present: true,
	}
}

// Empty creates a new empty Box
func Empty[T any]() Box[T] {
	return Box[T]{
		value:   nil,
		present: false,
	}
}
