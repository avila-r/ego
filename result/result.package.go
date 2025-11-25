package result

// Map applies a function to transform the value from T to R if Ok, otherwise returns the error
func Map[T, R any](r Result[T], f func(T) R) Result[R] {
	if r.IsError() {
		return Error[R](r.err)
	}
	return Ok(f(*r.value))
}

// MapOr applies a function to the value if Ok, otherwise returns the default value
func MapOr[T, R any](r Result[T], defaultValue R, f func(T) R) R {
	if r.IsError() {
		return defaultValue
	}
	return f(*r.value)
}

// MapOrElse applies a function to the value if Ok, otherwise calls the fallback function
func MapOrElse[T, R any](r Result[T], fallback func(error) R, f func(T) R) R {
	if r.IsError() {
		return fallback(r.err)
	}
	return f(*r.value)
}

// FlatMap applies a function that returns a Result and flattens it
func FlatMap[T, R any](r Result[T], f func(T) Result[R]) Result[R] {
	if r.IsError() {
		return Error[R](r.err)
	}
	return f(*r.value)
}

// AndThen chains operations that return Results (alias for FlatMap)
func AndThen[T, R any](r Result[T], f func(T) Result[R]) Result[R] {
	return FlatMap(r, f)
}

// Bind is another alias for FlatMap, common in functional programming
func Bind[T, R any](r Result[T], f func(T) Result[R]) Result[R] {
	return FlatMap(r, f)
}
