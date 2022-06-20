package generics

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrNotDone = errors.New("not done")

type Promising[T, U any] func(context.Context, T) (U, error)

type Waiter interface { // doesn't bring cocktails :P
	Wait() error
}

type Promise[T any] struct {
	result T
	err    error
	done   <-chan struct{}
}

func (p *Promise[T]) Resolve() (T, error) {
	<-p.done

	return p.result, p.err
}

func (p *Promise[T]) Try() (T, error) {
	select {
	case <-p.done:
		return p.result, p.err
	default:
		var zero T
		return zero, ErrNotDone
	}
}

func (p *Promise[T]) Wait() error {
	<-p.done
	return p.err
}

func Go[T, U any](ctx context.Context, t T, fn Promising[T, U]) *Promise[U] {
	done := make(chan struct{})
	result := Promise[U]{done: done}

	go func() {
		defer close(done)
		result.result, result.err = fn(ctx, t)
	}()

	return &result
}

func Wait(ws ...Waiter) error {
	var wg sync.WaitGroup
	wg.Add(len(ws))
	errChan := make(chan error, len(ws))
	done := make(chan struct{})

	for _, w := range ws {
		go func(w Waiter) {
			defer wg.Done()
			if err := w.Wait(); err != nil {
				errChan <- err
			}
		}(w)
	}

	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case err := <-errChan:
		return err
	case <-done: // nothing to do
	}

	return nil
}

func WithCancel[T, U any](fn Promising[T, U]) Promising[T, U] {
	return func(ctx context.Context, t T) (U, error) {
		var (
			val U
			err error
		)

		done := make(chan struct{})
		go func() {
			defer close(done)
			val, err = fn(ctx, t)
		}()

		select {
		case <-ctx.Done():
			return *new(U), ctx.Err()
		case <-done: // nothing to do
		}

		return val, err
	}
}

func Then[T, U any](ctx context.Context, p *Promise[T], fn Promising[T, U]) *Promise[U] {
	done := make(chan struct{})
	result := Promise[U]{done: done}

	go func() {
		defer close(done)
		first, err := p.Resolve()
		if err != nil {
			result.err = err
			return
		}

		second, err := fn(ctx, first)
		result.result = second
		result.err = err
	}()

	return &result
}

type debounce struct {
	mu        *sync.Mutex
	timer     *time.Timer
	callbacks []func()
	after     time.Duration
	done      bool
}

func (d *debounce) reset() *debounce {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.done {
		return d
	}

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.after, func() {
		for _, f := range d.callbacks {
			f()
		}
	})
	return d
}

func (d *debounce) cancel() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	d.done = true
}

// NewDebounce creates a debounced instance that delays invoking functions given until after wait milliseconds have elapsed.
func NewDebounce(delay time.Duration, fns ...func()) (func(), func()) {
	d := &debounce{
		after:     delay,
		mu:        new(sync.Mutex),
		callbacks: fns,
	}

	return func() {
		d.reset()
	}, d.cancel
}

// Invoke invokes a function N times until it returns valid output. Returning either the caught error or nil. When first argument is less than `1`, the function runs until a successful response is returned.
func Invoke(times int, fn func(int) error) (int, error) {
	var err error

	for i := 0; times <= 0 || i < times; i++ {
		err = fn(i)
		if err == nil {
			return i + 1, nil
		}
	}

	return times, err
}

// DelayedInvoke invokes a function N times until it returns valid output, with a pause between each call. Returning either the caught error or nil.
// When first argument is less than `1`, the function runs until a successful response is returned.
func DelayedInvoke(times int, delay time.Duration, fn func(int, time.Duration) error) (int, time.Duration, error) {
	var err error

	start := time.Now()

	for i := 0; times <= 0 || i < times; i++ {
		err = fn(i, time.Since(start))
		if err == nil {
			return i + 1, time.Since(start), nil
		}

		if times <= 0 || i+1 < times {
			time.Sleep(delay)
		}
	}

	return times, time.Since(start), err
}
