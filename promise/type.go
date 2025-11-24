package promise

// Concurrency determines whether operations run synchronously or asynchronously
type Concurrency int

const (
	Sync Concurrency = iota
	Async
)

func (c Concurrency) IsAsync() bool {
	return c == Async
}
