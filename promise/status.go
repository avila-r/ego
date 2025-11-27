package promise

// State represents the current state of a Promise
type State int

const (
	StateRunning State = iota
	StateCompleted
	StateFailed
	StateCancelled
)

var labels = []string{
	"Running",
	"Success",
	"Failed",
	"Cancelled",
}

func (s State) String() string {
	return labels[s]
}
