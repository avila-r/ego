package httpx

type Protocol struct {
	name  string
	major int
	minor int
}

var (
	ProtocolHTTP1_0 = Protocol{"HTTP", 1, 0}
	ProtocolHTTP1_1 = Protocol{"HTTP", 1, 1}
	ProtocolHTTP2   = Protocol{"HTTP", 2, 0}
)

func (p Protocol) Name() string {
	return p.name
}

func (p Protocol) Major() int {
	return p.major
}

func (p Protocol) Minor() int {
	return p.minor
}
