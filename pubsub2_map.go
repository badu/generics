package generics

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// entry is a *Topic[T] in the TopicsMap, but this cannot be made generic
// All uses of its fields must be atomic.
type entry struct {
	pointer unsafe.Pointer // *any
}

// TopicsMap is a concurrent read-mostly map, much like sync.Map with string keys.
type TopicsMap struct {
	value      atomic.Value // map[EventID]*entry
	mu         sync.Mutex   // mu must be held when using dirty or misses.
	dirtiesMap map[EventID]*entry
	misses     int
}

func newEntry(v any) *entry {
	return &entry{unsafe.Pointer(&v)}
}

func (e *entry) load() (any, bool) {
	if e == nil {
		return nil, false
	}

	p := atomic.LoadPointer(&e.pointer)
	if p == nil {
		// Nil means deleted.
		return nil, false
	}

	return *(*any)(p), true
}

func (e *entry) store(topic any) {
	atomic.StorePointer(&e.pointer, unsafe.Pointer(&topic))
}

func (e *entry) delete() *any {
	return (*any)(atomic.SwapPointer(&e.pointer, nil))
}

// Store sets the value at a key.
func (m *TopicsMap) Store(eventID EventID, topic any) {
	iMap, _ := m.value.Load().(map[EventID]*entry)
	e := iMap[eventID]
	if e != nil {
		e.store(topic)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// in case another goroutine set it while we were locking, reload
	iMap, _ = m.value.Load().(map[EventID]*entry)
	e = iMap[eventID]
	if e != nil {
		e.store(topic)
		return
	}

	e = m.dirtiesMap[eventID]
	if e == nil {
		m.miss() // making sure dirtiesMap is non-nil.
		m.dirtiesMap[eventID] = newEntry(topic)
		return
	}
	e.store(topic)
}

// Load gets the value at a key. ok is false if the key was not in the map.
func (m *TopicsMap) Load(eventID EventID) (any, bool) {
	mv, _ := m.value.Load().(map[EventID]*entry)
	e, ok := mv[eventID]
	if ok {
		return e.load()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// Reload e in case another goroutine set it while we were locking.
	mv, _ = m.value.Load().(map[EventID]*entry)
	e, ok = mv[eventID]
	if !ok {
		e, ok = m.dirtiesMap[eventID]
		m.miss()
	}
	return e.load()
}

// LoadOrStore gets the value at a key if it exists or stores and returns v if
// it does not. loaded is true if the value already existed.
func (m *TopicsMap) LoadOrStore(eventID EventID, topic any) (any, bool) {
	mv, _ := m.value.Load().(map[EventID]*entry)
	e, ok := mv[eventID]
	if ok {
		return e.load()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// Reload e in case another goroutine set it while we were locking.
	mv, _ = m.value.Load().(map[EventID]*entry)
	e, ok = mv[eventID]
	if ok {
		return e.load()
	}
	e, ok = m.dirtiesMap[eventID]
	// Whether we load or store, this is a miss.
	m.miss()
	if ok {
		return e.load()
	}
	m.dirtiesMap[eventID] = newEntry(topic)
	return topic, false
}

// Delete deletes the value at a key.
func (m *TopicsMap) Delete(eventID EventID) {
	m.LoadAndDelete(eventID)
}

// LoadAndDelete deletes the value at a key, returning its old value and whether it existed.
func (m *TopicsMap) LoadAndDelete(eventID EventID) (any, bool) {
	mv, _ := m.value.Load().(map[EventID]*entry)
	e := mv[eventID]
	if e != nil {
		if p := e.delete(); p != nil {
			return *p, true
		}
		return nil, false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// Reload e in case another goroutine set it while we were locking.
	mv, _ = m.value.Load().(map[EventID]*entry)
	e = mv[eventID]
	if e != nil {
		if p := e.delete(); p != nil {
			return *p, true
		}
		return nil, false
	}
	e = m.dirtiesMap[eventID]
	m.miss()
	if e != nil {
		if p := e.delete(); p != nil {
			return *p, true
		}
	}
	return nil, false
}

// Range calls f for each key and its corresponding value in the map. If f
// returns false, the iteration ceases. Note that Range is O(n) even if f
// returns false after a constant number of calls.
func (m *TopicsMap) Range(iterFn func(eventID EventID, value any) bool) {
	m.mu.Lock()
	// Force miss to promote.
	m.misses = len(m.dirtiesMap) - 1
	m.miss()
	mv, _ := m.value.Load().(map[EventID]*entry)
	m.mu.Unlock()

	for k, v := range mv {
		if r, ok := v.load(); ok {
			if !iterFn(k, r) {
				return
			}
		}
	}
}

// miss updates the miss counter and possibly promotes the dirty map. The caller must hold m.mu.
func (m *TopicsMap) miss() {
	m.misses++
	if m.misses < len(m.dirtiesMap) {
		return
	}
	dirtiesMap := m.dirtiesMap
	m.value.Store(dirtiesMap)
	m.dirtiesMap = make(map[EventID]*entry, len(dirtiesMap))
	for eventID, topic := range dirtiesMap {
		if atomic.LoadPointer(&topic.pointer) != nil {
			m.dirtiesMap[eventID] = topic
		}
	}
	m.misses = 0
}
