package generics

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

const (
	messagesPerTopic    = 200
	subscribersPerTopic = 200
)

func TestBus(t *testing.T) {
	t.Parallel()

	bus := NewPubSub[int](DefaultPublishChannelBufferSize)

	type result struct {
		key, val int
	}
	subscribersChan := make(chan result, subscribersPerTopic*messagesPerTopic)

	wg := sync.WaitGroup{}
	wg.Add(subscribersPerTopic)
	for i := 0; i < subscribersPerTopic; i++ {
		go func(i int) {
			defer wg.Done()
			err := bus.Listen(DefaultSubscribeChannelBufferSize, func(message int) {
				subscribersChan <- result{
					key: i,
					val: message,
				}
			})
			if err != nil {
				t.Errorf("failed to subscribe the topic: %s", err)
				return
			}
		}(i)
	}
	wg.Wait()

	wg.Add(messagesPerTopic)
	for i := 0; i < messagesPerTopic; i++ {
		go func(i int) {
			defer wg.Done()
			if err := bus.Dispatch(i); err != nil {
				t.Errorf("failed to publish a message to the topic: %s", err)
				return
			}
		}(i)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		if err := bus.Start(ctx); err != nil {
			t.Errorf("the topic has aborted: %s", err)
			return
		}

		done <- struct{}{}
	}()
	wg.Wait()

	cancel()
	<-done

	close(subscribersChan)

	seen := make(map[int]map[int]struct{})
	for subscriber := range subscribersChan {
		_, ok := seen[subscriber.key]
		if !ok {
			seen[subscriber.key] = make(map[int]struct{})
		}

		seen[subscriber.key][subscriber.val] = struct{}{}
	}

	found := make(map[int]map[int]struct{})
	for i := 0; i < subscribersPerTopic; i++ {
		found[i] = make(map[int]struct{})
		for j := 0; j < messagesPerTopic; j++ {
			found[i][j] = struct{}{}
		}
	}

	if !reflect.DeepEqual(seen, found) {
		t.Error("seen and found should be equal")
	}
}

func TestBusErrors(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	bus := NewPubSub[struct{}](DefaultPublishChannelBufferSize)

	done := make(chan struct{})
	go func() {
		if err := bus.Start(ctx); err != nil {
			t.Errorf("the bus has aborted: %s", err)
			return
		}

		done <- struct{}{}
	}()

	cancel()
	<-done

	err := bus.Dispatch(struct{}{})
	if err != ErrNotRunning {
		t.Errorf("unexpected error for the Dispatch: %s", err)
	}

	err = bus.Listen(DefaultSubscribeChannelBufferSize, func(message struct{}) {})
	if err != ErrNotRunning {
		t.Errorf("unexpected error for the Listen: %s", err)
	}

	err = bus.Start(ctx)
	if err != ErrAlreadyStarted {
		t.Errorf("unexpected error for the Listen: %s", err)
	}
}

// =============
// example below
// =============

type greetingMessage struct {
	greeting string
}

func Example() {

	bus := NewPubSub[greetingMessage](DefaultPublishChannelBufferSize) // the topic with a type which you want to publish and subscribe.

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		if err := bus.Start(ctx); err != nil { // Start the topic. This call of Start blocks until the context is canceled.
			println(err)
			return
		}

		done <- struct{}{}
	}()

	if err := bus.Dispatch(greetingMessage{greeting: "Hello, badu!"}); err != nil { // Dispatch a message to the topic. This call of Dispatch is non-blocking.
		println(err)
		return
	}

	// Listen the topic. This call of Listen is non-blocking.
	err := bus.Listen(DefaultSubscribeChannelBufferSize, func(message greetingMessage) {
		println(message.greeting)
	})
	if err != nil {
		println(err)
		return
	}

	cancel()
	<-done
}
