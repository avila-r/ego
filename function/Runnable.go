package function

type Runnable interface {
	Run()
}

type DefaultRunnable struct {
	runnable func()
}

func (d *DefaultRunnable) Run() {
	d.runnable()
}

func NewRunnable(runnable func()) Runnable {
	return &DefaultRunnable{runnable: runnable}
}
