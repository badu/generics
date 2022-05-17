package generics

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrNotDone = errors.New("not done")

type Promising[T, V any] func(context.Context, T) (V, error)

type Waiter interface { // doesn't bring cocktails :P
	Wait() error
}

type Promise[V any] struct {
	result V
	err    error
	done   <-chan struct{}
}

func (p *Promise[V]) Resolve() (V, error) {
	<-p.done
	return p.result, p.err
}

func (p *Promise[V]) AttemptResolve() (V, error) {
	select {
	case <-p.done:
		return p.result, p.err
	default:
		var zero V
		return zero, ErrNotDone
	}
}

func Go[T, V any](ctx context.Context, t T, fn Promising[T, V]) *Promise[V] {
	done := make(chan struct{})
	p := Promise[V]{
		done: done,
	}
	go func() {
		defer close(done)
		p.result, p.err = fn(ctx, t)
	}()
	return &p
}

func (p *Promise[V]) Wait() error {
	<-p.done
	return p.err
}

func Wait(ws ...Waiter) error {
	var wg sync.WaitGroup
	wg.Add(len(ws))
	errChan := make(chan error, len(ws))
	done := make(chan struct{})
	for _, w := range ws {
		go func(w Waiter) {
			defer wg.Done()
			err := w.Wait()
			if err != nil {
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
	case <-done:
	}
	return nil
}

func WithCancel[T, V any](fn Promising[T, V]) Promising[T, V] {
	return func(ctx context.Context, t T) (V, error) {
		var (
			val V
			err error
		)

		done := make(chan struct{})
		go func() {
			defer close(done)
			val, err = fn(ctx, t)
		}()

		select {
		case <-ctx.Done():
			return *new(V), ctx.Err()
		case <-done: // nothing to do
		}

		return val, err
	}
}

func Then[T, V any](ctx context.Context, p *Promise[T], fn Promising[T, V]) *Promise[V] {
	done := make(chan struct{})
	out := Promise[V]{
		done: done,
	}
	go func() {
		defer close(done)
		val, err := p.Resolve()
		if err != nil {
			out.err = err
			return
		}
		val2, err := fn(ctx, val)
		out.result = val2
		out.err = err
	}()
	return &out
}

type debounce struct {
	after     time.Duration
	mu        *sync.Mutex
	timer     *time.Timer
	done      bool
	callbacks []func()
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

// Attempt invokes a function N times until it returns valid output. Returning either the caught error or nil. When first argument is less than `1`, the function runs until a successful response is returned.
func Attempt(times int, fn func(int) error) (int, error) {
	var err error

	for i := 0; times <= 0 || i < times; i++ {
		err = fn(i)
		if err == nil {
			return i + 1, nil
		}
	}

	return times, err
}

// AttemptAfter invokes a function N times until it returns valid output, with a pause between each call. Returning either the caught error or nil.
// When first argument is less than `1`, the function runs until a successful response is returned.
func AttemptAfter(times int, delay time.Duration, fn func(int, time.Duration) error) (int, time.Duration, error) {
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
