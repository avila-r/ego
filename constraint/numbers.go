package constraint

type (
	Complex interface {
		~complex64 | ~complex128
	}

	Float interface {
		~float32 | ~float64
	}

	Unsigned interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}

	Signed interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}

	Integer interface {
		Signed | Unsigned
	}

	Arithmetic interface {
		Integer | Float | Complex
	}
)
