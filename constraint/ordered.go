package constraint

type Ordered interface {
	Comparable | Unsigned
}
