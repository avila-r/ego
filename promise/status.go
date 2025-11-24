package promise

// State represents the current state of a Promise
type State int

const (
	StateRunning State = iota
	StateSuccess
	StateFailed
	StateCancelled
)

func (s State) String() string {
	return []string{"Running", "Success", "Failed", "Cancelled"}[s]
}
