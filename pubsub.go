package generics

import (
	"context"
	"errors"
	"sync"
)

const (
	DefaultPublishChannelBufferSize   = 100
	DefaultSubscribeChannelBufferSize = 100
)

var (
	ErrAlreadyStarted = errors.New("this topic has been already started")
	ErrNotRunning     = errors.New("this topic has been already closed")
)

type runningMutex struct {
	sync.RWMutex
	running bool
}

// Topic is a publish-subscribe topic.
type Topic[T any] struct {
	publish          chan T
	done             chan struct{}
	wg               sync.WaitGroup
	runningMu        runningMutex
	endingMu         runningMutex
	subscribersStore struct {
		sync.RWMutex
		subscribers []*subscription[T]
	}
}

// NewPubSub creates a new topic.
func NewPubSub[T any](bufferSize int) *Topic[T] {
	result := &Topic[T]{
		publish: make(chan T, bufferSize),
		done:    make(chan struct{}),
	}
	return result
}

// Start starts the topic and blocks until the context is canceled.
// When the passed context is canceled, Start waits for all published messages to be processed by all subscribers.
// Once Start is called, then Start returns ErrAlreadyStarted.
func (t *Topic[T]) Start(ctx context.Context) error {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	if t.runningMu.running {
		return ErrAlreadyStarted
	}
	t.runningMu.running = true

	go func() {
		defer close(t.done)

		for {
			message, ok := <-t.publish
			if !ok {
				return
			}

			go func(message T) {
				defer t.wg.Done()

				var wg sync.WaitGroup

				t.subscribersStore.RLock()
				for _, subscriber := range t.subscribersStore.subscribers {
					wg.Add(1)
					go func(s *subscription[T]) {
						defer wg.Done()
						s.message <- message
					}(subscriber)
				}
				t.subscribersStore.RUnlock()

				wg.Wait()
			}(message)
		}
	}()

	<-ctx.Done()
	t.stop()

	return nil
}

// Dispatch publishes a message to the topic. This method is non-blocking and concurrent-safe. Returns ErrNotRunning if the topic has been already closed.
func (t *Topic[T]) Dispatch(message T) error {
	t.endingMu.RLock()
	defer t.endingMu.RUnlock() // avoid the race condition between Topic.Dispatch and Topic.stop, defer is necessary.
	if t.endingMu.running {
		return ErrNotRunning
	}

	t.wg.Add(1)
	t.publish <- message

	return nil
}

// Listen registers the passed function as a subscriber to the topic. This method is non-blocking and concurrent-safe. Function passed to Listen is called when a message is published to the topic, obviously. Returns ErrNotRunning if the topic has been already closed.
func (t *Topic[T]) Listen(bufferSize int, subscriber func(message T)) error {
	t.endingMu.RLock()
	if t.endingMu.running {
		return ErrNotRunning
	}
	t.endingMu.RUnlock()

	newSubscriber := &subscription[T]{
		message:    make(chan T, bufferSize),
		done:       make(chan struct{}),
		subscriber: subscriber,
	}

	t.subscribersStore.Lock()
	t.subscribersStore.subscribers = append(t.subscribersStore.subscribers, newSubscriber)
	t.subscribersStore.Unlock()

	go func() {
		defer close(newSubscriber.done)

		for {
			message, ok := <-newSubscriber.message
			if !ok {
				newSubscriber.wg.Wait()
				return
			}

			newSubscriber.wg.Add(1)

			go func() {
				defer newSubscriber.wg.Done()
				subscriber(message)
			}()
		}
	}()

	return nil
}

type subscription[T any] struct {
	message    chan T
	done       chan struct{}
	wg         sync.WaitGroup
	subscriber func(message T)
}

func (t *Topic[T]) stop() {
	t.endingMu.Lock()
	t.endingMu.running = true
	t.endingMu.Unlock()

	t.wg.Wait()

	close(t.publish)

	var wg sync.WaitGroup

	t.subscribersStore.RLock()
	for _, subscriber := range t.subscribersStore.subscribers {
		close(subscriber.message)

		wg.Add(1)
		go func(s *subscription[T]) {
			defer wg.Done()
			<-s.done
		}(subscriber)
	}
	t.subscribersStore.RUnlock()

	wg.Wait()
	<-t.done
}
