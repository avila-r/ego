package promise

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Promise represents an asynchronous computation that will eventually produce a value or error
type Promise[T any] struct {
	mu     sync.RWMutex
	done   chan Void
	once   sync.Once
	value  T
	err    error
	state  State
	ctx    context.Context
	cancel context.CancelFunc
}

type (
	Void struct{}
)

type (
	Runnable func()

	Supplier[T any] func() (T, error)

	Pipeline[T any] func(T) T

	Consumer[T any] func(T)

	Mapper[T any, R any] func(T) R

	Composer[T any, R any] func(T) *Promise[R]

	Recovery[T any] func(error) T

	Handler[T any, R any] func(T, error) R
)

// Run creates a Promise that executes a function with no return value
func Run(function Runnable) *Promise[Void] {
	return Of(func() (Void, error) {
		function()
		return Void{}, nil
	})
}

// Empty creates an incomplete Promise that can be manually completed
func Empty[T any]() *Promise[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return &Promise[T]{
		done:   make(chan Void),
		state:  StateRunning,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Of creates a Promise from a supplier function
func Of[T any](supplier Supplier[T]) *Promise[T] {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Promise[T]{
		done:   make(chan Void),
		state:  StateRunning,
		ctx:    ctx,
		cancel: cancel,
	}

	go func() {
		defer p.close()
		value, err := supplier()
		p.mu.Lock()
		defer p.mu.Unlock()
		p.value = value
		p.err = err
		if err != nil {
			p.state = StateFailed
		} else {
			p.state = StateSuccess
		}
	}()

	return p
}

// Supply is an alias for Of
func Supply[T any](supplier Supplier[T]) *Promise[T] {
	return Of(supplier)
}

// Completed creates a Promise that is already completed with a value
func Completed[T any](value T) *Promise[T] {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Promise[T]{
		done:   make(chan Void),
		value:  value,
		state:  StateSuccess,
		ctx:    ctx,
		cancel: cancel,
	}
	p.close()
	return p
}

// Done is an alias for Completed
func Done[T any](value T) *Promise[T] {
	return Completed(value)
}

// Then applies a transformation function and returns the same Promise
func (p *Promise[T]) Then(pipeline Pipeline[T]) *Promise[T] {
	return p.ThenWithConcurrency(Async, pipeline)
}

// ThenWithConcurrency applies a pipeline with specified concurrency
func (p *Promise[T]) ThenWithConcurrency(concurrency Concurrency, pipeline Pipeline[T]) *Promise[T] {
	<-p.done
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.err != nil {
		return p
	}

	if concurrency.IsAsync() {
		return Of(func() (T, error) {
			return pipeline(p.value), nil
		})
	}

	p.value = pipeline(p.value)
	return p
}

// ThenAccept consumes the value and returns a Void Promise
func (p *Promise[T]) ThenAccept(consumer Consumer[T]) *Promise[Void] {
	return p.ThenAcceptWithConcurrency(Async, consumer)
}

// ThenAcceptWithConcurrency consumes the value with specified concurrency
func (p *Promise[T]) ThenAcceptWithConcurrency(concurrency Concurrency, consumer Consumer[T]) *Promise[Void] {
	if concurrency.IsAsync() {
		return Of(func() (Void, error) {
			<-p.done
			p.mu.RLock()
			defer p.mu.RUnlock()
			if p.err != nil {
				return Void{}, p.err
			}
			consumer(p.value)
			return Void{}, nil
		})
	}

	<-p.done
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.err != nil {
		return Completed(Void{})
	}
	consumer(p.value)
	return Completed(Void{})
}

// ThenRun executes an action and returns a Void Promise
func (p *Promise[T]) ThenRun(action Runnable) *Promise[Void] {
	return p.ThenRunWithConcurrency(Async, action)
}

// ThenRunWithConcurrency executes an action with specified concurrency
func (p *Promise[T]) ThenRunWithConcurrency(concurrency Concurrency, action Runnable) *Promise[Void] {
	if concurrency.IsAsync() {
		return Of(func() (Void, error) {
			<-p.done
			action()
			return Void{}, nil
		})
	}

	<-p.done
	action()
	return Completed(Void{})
}

// Map transforms the Promise value to a different type
func Map[T any, R any](p *Promise[T], mapper Mapper[T, R]) *Promise[R] {
	return MapWithConcurrency(p, Async, mapper)
}

// MapWithConcurrency transforms the Promise value with specified concurrency
func MapWithConcurrency[T any, R any](p *Promise[T], concurrency Concurrency, mapper Mapper[T, R]) *Promise[R] {
	if concurrency.IsAsync() {
		return Of(func() (R, error) {
			<-p.done
			p.mu.RLock()
			defer p.mu.RUnlock()
			if p.err != nil {
				var zero R
				return zero, p.err
			}
			return mapper(p.value), nil
		})
	}

	<-p.done
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.err != nil {
		var zero R
		return Completed(zero)
	}
	return Completed(mapper(p.value))
}

// Compose chains Promises together (flatMap)
func Compose[T any, R any](p *Promise[T], composer Composer[T, R]) *Promise[R] {
	return ComposeWithConcurrency(p, Async, composer)
}

// ComposeWithConcurrency chains Promises with specified concurrency
func ComposeWithConcurrency[T any, R any](p *Promise[T], concurrency Concurrency, composer Composer[T, R]) *Promise[R] {
	if concurrency.IsAsync() {
		return Of(func() (R, error) {
			<-p.done
			p.mu.RLock()
			if p.err != nil {
				p.mu.RUnlock()
				var zero R
				return zero, p.err
			}
			val := p.value
			p.mu.RUnlock()

			nextPromise := composer(val)
			<-nextPromise.done
			nextPromise.mu.RLock()
			defer nextPromise.mu.RUnlock()
			return nextPromise.value, nextPromise.err
		})
	}

	<-p.done
	p.mu.RLock()
	if p.err != nil {
		p.mu.RUnlock()
		var zero R
		return Completed(zero)
	}
	val := p.value
	p.mu.RUnlock()

	return composer(val)
}

// Exceptionally handles errors and provides a recovery value
func (p *Promise[T]) Exceptionally(recovery Recovery[T]) *Promise[T] {
	return p.ExceptionallyWithConcurrency(Async, recovery)
}

// Recover is an alias for Exceptionally
func (p *Promise[T]) Recover(recovery Recovery[T]) *Promise[T] {
	return p.Exceptionally(recovery)
}

// ExceptionallyWithConcurrency handles errors with specified concurrency
func (p *Promise[T]) ExceptionallyWithConcurrency(concurrency Concurrency, recovery Recovery[T]) *Promise[T] {
	if concurrency.IsAsync() {
		return Of(func() (T, error) {
			<-p.done
			p.mu.RLock()
			defer p.mu.RUnlock()
			if p.err != nil {
				return recovery(p.err), nil
			}
			return p.value, nil
		})
	}

	<-p.done
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.err != nil {
		return Completed(recovery(p.err))
	}
	return Completed(p.value)
}

// Handle processes both success and error cases
func Handle[T any, R any](p *Promise[T], handler Handler[T, R]) *Promise[R] {
	return HandleWithConcurrency(p, Async, handler)
}

// HandleWithConcurrency processes both cases with specified concurrency
func HandleWithConcurrency[T any, R any](p *Promise[T], concurrency Concurrency, handler Handler[T, R]) *Promise[R] {
	if concurrency.IsAsync() {
		return Of(func() (R, error) {
			<-p.done
			p.mu.RLock()
			defer p.mu.RUnlock()
			return handler(p.value, p.err), nil
		})
	}

	<-p.done
	p.mu.RLock()
	defer p.mu.RUnlock()
	return Completed(handler(p.value, p.err))
}

// Join blocks until the Promise completes and returns the value or panics on error
func (p *Promise[T]) Join() T {
	<-p.done
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.err != nil {
		panic(p.err)
	}
	return p.value
}

// Get blocks until the Promise completes and returns the value and error
func (p *Promise[T]) Get() (T, error) {
	<-p.done
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.value, p.err
}

// Complete manually completes the Promise with a value
func (p *Promise[T]) Complete(value T) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != StateRunning {
		return false
	}

	p.value = value
	p.state = StateSuccess
	p.close()
	return true
}

// CompleteExceptionally manually completes the Promise with an error
func (p *Promise[T]) CompleteExceptionally(err error) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != StateRunning {
		return false
	}

	p.err = err
	p.state = StateFailed
	p.close()
	return true
}

// Timeout sets a timeout for the Promise
func (p *Promise[T]) Timeout(duration time.Duration) *Promise[T] {
	go func() {
		select {
		case <-p.done:
			return
		case <-time.After(duration):
			p.CompleteExceptionally(errors.New("promise timeout"))
		}
	}()
	return p
}

// Threshold is an alias for Timeout
func (p *Promise[T]) Threshold(duration time.Duration) *Promise[T] {
	return p.Timeout(duration)
}

// CompleteOnTimeout completes with a default value if timeout occurs
func (p *Promise[T]) CompleteOnTimeout(value T, duration time.Duration) *Promise[T] {
	go func() {
		select {
		case <-p.done:
			return
		case <-time.After(duration):
			p.Complete(value)
		}
	}()
	return p
}

// IsDone returns true if the Promise has completed
func (p *Promise[T]) IsDone() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state != StateRunning
}

// IsCompleted is an alias for IsDone
func (p *Promise[T]) IsCompleted() bool {
	return p.IsDone()
}

// IsSuccess returns true if the Promise completed successfully
func (p *Promise[T]) IsSuccess() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state == StateSuccess
}

// IsCancelled returns true if the Promise was cancelled
func (p *Promise[T]) IsCancelled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state == StateCancelled
}

// IsError returns true if the Promise completed with an error
func (p *Promise[T]) IsError() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state == StateFailed
}

// IsCompletedExceptionally is an alias for IsError
func (p *Promise[T]) IsCompletedExceptionally() bool {
	return p.IsError()
}

// State returns the current state of the Promise
func (p *Promise[T]) State() State {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// Cancel cancels the Promise's context
func (p *Promise[T]) Cancel() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.state == StateRunning {
		p.state = StateCancelled
		p.cancel()
		p.close()
	}
}

// Context returns the Promise's context
func (p *Promise[T]) Context() context.Context {
	return p.ctx
}

// close safely closes the done channel only once
func (p *Promise[T]) close() {
	p.once.Do(func() {
		close(p.done)
	})
}
