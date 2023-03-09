package generics

import (
	"sync"
	"sync/atomic"
)

type EventID string

type Event interface {
	EventID() EventID
}

type Topic[T Event] struct {
	subscriptions []*Subscription[T]
	rwMu          sync.RWMutex
	pool          sync.Pool
}

func NewTopic[T Event]() *Topic[T] {
	result := &Topic[T]{}
	result.pool.New = func() any {
		return &Subscription[T]{
			topic:    result,
			callback: nil,
		}
	}
	return result
}

type Subscription[T Event] struct {
	topic    *Topic[T]
	callback func(event T)
}

func (t *Topic[T]) Subscribe(callback func(v T)) *Subscription[T] {
	result := t.pool.Get().(*Subscription[T])
	result.callback = callback
	result.topic = t

	t.rwMu.Lock()
	t.subscriptions = append(t.subscriptions, result)
	t.rwMu.Unlock()

	return result
}

func (t *Topic[T]) Unsubscribe(subscription *Subscription[T]) {
	t.rwMu.Lock()
	for i := range t.subscriptions {
		if t.subscriptions[i] != subscription {
			continue
		}

		t.subscriptions[i] = t.subscriptions[len(t.subscriptions)-1]
		t.subscriptions[len(t.subscriptions)-1] = nil
		t.subscriptions = t.subscriptions[:len(t.subscriptions)-1]
		break
	}
	t.rwMu.Unlock()

	subscription.callback = nil
	t.pool.Put(subscription)
}

func (t *Topic[T]) NumSubscribers() uint64 {
	t.rwMu.RLock()
	result := uint64(len(t.subscriptions))
	t.rwMu.RUnlock()
	return result
}

func (s *Subscription[T]) Unsubscribe() {
	s.topic.Unsubscribe(s)
}

func (s *Subscription[T]) Topic() *Topic[T] {
	return s.topic
}

// TODO : add param useGoroutines bool, to spin the callback on a goroutine
func (t *Topic[T]) Broadcast(event T) {
	t.rwMu.RLock()
	for topic := range t.subscriptions {
		t.subscriptions[topic].callback(event)
	}
	t.rwMu.RUnlock()
}

type ChannelTopic[T Event] struct {
	subscription *Subscription[T]
	name         EventID
	OnMessage    func(event T)
}

func NewBroadcastChannel[T Event](channelsMap *TopicsMap) *ChannelTopic[T] {
	var event T
	topic, ok := channelsMap.Load(event.EventID())
	if topic != ok {
		topic, _ = channelsMap.LoadOrStore(event.EventID(), NewTopic[T]())
	}

	channelTopic := &ChannelTopic[T]{
		name: event.EventID(),
	}

	listener := topic.(*Topic[T]).Subscribe(func(v T) {
		o := channelTopic.OnMessage
		if o != nil {
			o(v)
		}
	})
	channelTopic.subscription = listener

	return channelTopic
}

func (ct ChannelTopic[T]) Name() string {
	return string(ct.name)
}

func (ct ChannelTopic[T]) Broadcast(event T) {
	ct.subscription.Topic().Broadcast(event)
}

func (ct ChannelTopic[T]) Close() {
	ct.subscription.Unsubscribe()
}

func PublishEvent[T Event](busMap *TopicsMap, event T) {
	topic, ok := busMap.Load(event.EventID())
	if !ok || topic == nil {
		topic, _ = busMap.LoadOrStore(event.EventID(), NewTopic[T]())
	}
	topic.(*Topic[T]).Broadcast(event)
}

type BusSubscription[T Event] struct {
	topic        EventID
	subscription *Subscription[T]
	stop         atomic.Uint32
}

func SubscribeEvent[T Event](busMap *TopicsMap, callback func(event T) bool) *BusSubscription[T] {
	var event T
	topic, ok := busMap.Load(event.EventID())
	if !ok || topic == nil {
		topic, _ = busMap.LoadOrStore(event.EventID(), NewTopic[T]())
	}
	var result BusSubscription[T]
	result.topic = event.EventID()

	result.subscription = topic.(*Topic[T]).Subscribe(func(v T) {
		if result.stop.Load() == 1 {
			return
		}

		unsub := callback(v)
		if unsub {
			result.Unsubscribe()
		}

	})

	return &result
}

func (bs *BusSubscription[T]) Unsubscribe() {
	if bs.stop.CompareAndSwap(0, 1) {
		go bs.subscription.Unsubscribe()
	}
}

func (bs *BusSubscription[T]) String() string {
	return "topic `" + string(bs.topic) + "`"
}
