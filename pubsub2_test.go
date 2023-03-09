package generics_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/badu/generics"
)

const (
	Uint32EventType         generics.EventID = "Uint32Event"
	EmptyStructEventType    generics.EventID = "EmptyStructEvent"
	UserRegisteredEventType generics.EventID = "UserRegisteredEvent"
	SMSRequestEventType     generics.EventID = "SMSRequestEvent"
	SMSSentEventType        generics.EventID = "SmsSentEvent"
)

type UserRegisteredEvent struct {
	UserName string
}

func (e UserRegisteredEvent) EventID() generics.EventID {
	return UserRegisteredEventType
}

type SMSRequestEvent struct {
	Number  string
	Message string
}

func (e SMSRequestEvent) EventID() generics.EventID {
	return SMSRequestEventType
}

type SMSSentEvent struct {
	Request SMSRequestEvent
	Status  string
}

func (e SMSSentEvent) EventID() generics.EventID {
	return SMSSentEventType
}

type Uint32Event struct {
	u uint32
}

func (e Uint32Event) EventID() generics.EventID {
	return Uint32EventType
}

type EmptyStructEvent struct{}

func (e EmptyStructEvent) EventID() generics.EventID {
	return EmptyStructEventType
}

func TestBroadcastChannel(t *testing.T) {
	var channelsMap = generics.TopicsMap{}
	bc := generics.NewBroadcastChannel[Uint32Event](&channelsMap)
	defer bc.Close()

	var ctr uint32
	bc.OnMessage = func(v Uint32Event) {
		ctr += v.u
	}
	for i := 0; i < 100; i++ {
		bc.Broadcast(Uint32Event{u: 100})
	}
	if ctr != 10000 {
		t.Errorf("ctr == %d, want 10000", ctr)
	}
}

func TestEventBus(t *testing.T) {
	var ctr atomic.Uint64
	var busMap generics.TopicsMap
	for i := 0; i < 8; i++ {
		ss := generics.SubscribeEvent[EmptyStructEvent](&busMap, func(EmptyStructEvent) bool {
			ctr.Add(1)
			return false
		})
		defer ss.Unsubscribe()
	}
	for i := 0; i < 100; i++ {
		generics.PublishEvent[EmptyStructEvent](&busMap, EmptyStructEvent{})
	}

	if ctr.Load() != 800 {
		t.Errorf("ctr=%d, ctr != 800", ctr.Load())
	}
}

type TestingB struct {
	*testing.B
}

func (e TestingB) EventID() generics.EventID {
	return "TestingBEvent"
}

func TestTopicBroadcast(t *testing.T) {
	tt := generics.NewTopic[Uint32Event]()
	var ctr uint32
	s := tt.Subscribe(func(v Uint32Event) {
		atomic.AddUint32(&ctr, v.u)
	})
	for i := 0; i < 100; i++ {
		s.Topic().Broadcast(Uint32Event{u: 100})
	}
	s.Unsubscribe()
	for i := 0; i < 100; i++ {
		s.Topic().Broadcast(Uint32Event{u: 100})
	}

	if ctr != 10000 {
		t.Fatalf("expected ctr == 10000, got %d", ctr)
	}
}

func BenchmarkBroadcast_0008(b *testing.B) {
	tt := generics.NewTopic[Uint32Event]()
	var ctr uint32
	for i := 0; i < 8; i++ {
		tt.Subscribe(func(v Uint32Event) {
			atomic.AddUint32(&ctr, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tt.Broadcast(Uint32Event{u: 1})
		}
	})
}

func BenchmarkBroadcast_0256(b *testing.B) {
	tt := generics.NewTopic[Uint32Event]()
	var ctr uint32
	for i := 0; i < 256; i++ {
		tt.Subscribe(func(v Uint32Event) {
			atomic.AddUint32(&ctr, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tt.Broadcast(Uint32Event{u: 1})
		}
	})
}

func BenchmarkBroadcast_1k(b *testing.B) {
	tt := generics.NewTopic[Uint32Event]()
	var ctr uint32
	for i := 0; i < 1024; i++ {
		tt.Subscribe(func(v Uint32Event) {
			atomic.AddUint32(&ctr, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tt.Broadcast(Uint32Event{u: 1})
		}
	})
}

func BenchmarkBroadcast_2k(b *testing.B) {
	tt := generics.NewTopic[Uint32Event]()
	var ctr uint32
	for i := 0; i < 2048; i++ {
		tt.Subscribe(func(v Uint32Event) {
			atomic.AddUint32(&ctr, v.u)
		})
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			tt.Broadcast(Uint32Event{u: 1})
		}
	})
}

func BenchmarkEventBus(b *testing.B) {
	var busMap generics.TopicsMap
	for i := 0; i < 8; i++ {
		ss := generics.SubscribeEvent[TestingB](&busMap, func(v TestingB) bool {
			return false
		})
		defer ss.Unsubscribe()
	}

	b.RunParallel(
		func(p *testing.PB) {
			for p.Next() {
				generics.PublishEvent[TestingB](&busMap, TestingB{b})
			}
		},
	)
}

type UserSvc struct {
	busMap *generics.TopicsMap
	t      *testing.T
}

func NewUserService(busMap *generics.TopicsMap, t *testing.T) UserSvc {
	result := UserSvc{busMap: busMap, t: t}
	if busMap != nil {
		generics.SubscribeEvent[SMSSentEvent](busMap, result.Handler)
	}
	return result
}

func (us *UserSvc) Handler(e SMSSentEvent) bool {
	us.t.Logf("SMS was sent (we know nothing about email though)")
	return false
}

func (us *UserSvc) RegisterUser(name, phoneNumber string) {
	if us.busMap != nil {
		generics.PublishEvent[UserRegisteredEvent](us.busMap, UserRegisteredEvent{UserName: name})
		generics.PublishEvent[SMSRequestEvent](us.busMap, SMSRequestEvent{Message: fmt.Sprintf("welcome %s", name), Number: phoneNumber})
	}
}

type EmailSvc struct {
	busMap *generics.TopicsMap
	t      *testing.T
}

func NewEmailService(busMap *generics.TopicsMap, t *testing.T) EmailSvc {
	result := EmailSvc{busMap: busMap, t: t}
	if busMap != nil {
		generics.SubscribeEvent[UserRegisteredEvent](busMap, result.Handler)
	}
	return result
}

func (es *EmailSvc) Handler(e UserRegisteredEvent) bool {
	es.t.Logf("Sending welcome email to %q", e.UserName)
	return false
}

type SMSSvc struct {
	busMap *generics.TopicsMap
	t      *testing.T
}

func NewSMSService(busMap *generics.TopicsMap, t *testing.T) SMSSvc {
	result := SMSSvc{busMap: busMap, t: t}
	if busMap != nil {
		generics.SubscribeEvent[SMSRequestEvent](busMap, result.Handler)
	}
	return result
}

func (ss *SMSSvc) Handler(e SMSRequestEvent) bool {
	if ss.busMap == nil {
		ss.t.Logf("skipping. bus map is nil")
		return true
	}
	ss.t.Logf("SMS Sent to number %s : %q", e.Number, e.Message)
	<-time.After(1 * time.Second) // simulate heavy op
	ss.t.Logf("Replying that SMS was sent")
	generics.PublishEvent[SMSSentEvent](ss.busMap, SMSSentEvent{
		Request: e,
		Status:  "sent",
	})
	return false
}

func TestManyListenersManySubscribers(t *testing.T) {
	var busMap generics.TopicsMap
	uSvc := NewUserService(&busMap, t)
	NewEmailService(&busMap, t)
	NewSMSService(&busMap, t)
	uSvc.RegisterUser("Badu", "0742.222.222")
	<-time.After(3 * time.Second)
}
