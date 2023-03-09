package generics_test

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/badu/generics"
)

func TestMap(t *testing.T) {
	var m generics.TopicsMap
	keys := strings.Fields("a b c d e f g h i j k l m n o p q r s t u v w x y z")
	ch := make(chan error, len(keys))
	start := make(chan bool)
	stop := make(chan bool)
	for _, k := range keys {
		go func(k string) {
			<-start
			p, ok := m.Load(generics.EventID(k))
			if ok {
				ch <- fmt.Errorf("unexpected successful initial load of key %q", k)
			}
			if p != nil {
				ch <- fmt.Errorf("initial load of key %q produced nonzero value %#v", k, p)
			}
			m.Store(generics.EventID(k), -1)
			for i := 0; i < 1e4; i++ {
				p, ok := m.Load(generics.EventID(k))
				if !ok {
					ch <- fmt.Errorf("load %d of key %q failed", i, k)
				}
				if p != i-1 {
					ch <- fmt.Errorf("load %d of key %q gave wrong result: want %d, got %v", i, k, i-1, p)
				}
				m.Store(generics.EventID(k), i)
			}
			ch <- nil
		}(k)
	}
	go func() {
		<-start
		for {
			select {
			case <-stop:
				return
			default:
				m.Load("nonexistent key to trigger frequent dirty map promotions")
			}
		}
	}()
	close(start)
	for i := 0; i < len(keys); i++ {
		err := <-ch
		if err != nil {
			t.Error(err)
			i-- // don't count errors, only nils
		}
	}
	close(stop)
}

func TestStoreOnly(t *testing.T) {
	var m generics.TopicsMap
	var wg sync.WaitGroup
	ch := make(chan bool)
	n := runtime.GOMAXPROCS(0)
	wg.Add(n)
	f := func(v int) {
		k := string(rune(v + '0'))
		<-ch
		for i := 0; i < 1e5; i++ {
			m.Store(generics.EventID(k), v)
		}
		wg.Done()
	}
	for i := 0; i < n; i++ {
		go f(i)
	}
	close(ch)
	wg.Wait()
	// No correct result, except that the race detector shouldn't complain.
}

func TestMapKey(t *testing.T) {
	var m generics.TopicsMap
	var wg sync.WaitGroup
	n := runtime.GOMAXPROCS(0)
	wg.Add(n)
	if _, ok := m.Load("key"); ok {
		t.Fatal("key exists in zero value")
	}
	ch := make(chan bool)
	f := func(v int) {
		<-ch
		for i := 0; i < 1e5; i++ {
			m.Load("key")
			m.Store("key", v)
		}
		wg.Done()
	}
	for i := 0; i < n; i++ {
		go f(i)
	}
	close(ch)
	wg.Wait()
	// No correct result, except that the race detector shouldn't complain.
}

func TestLoadOrStore(t *testing.T) {
	var m generics.TopicsMap
	n := runtime.GOMAXPROCS(0)
	ch := make(chan bool)
	errs := make(chan error, n)
	var stores int32
	stored := int32(-1)
	f := func(v int32) {
		<-ch
		for i := 0; i < 1e5; i++ {
			r, ok := m.LoadOrStore("key", v)
			if ok {
				if s := atomic.LoadInt32(&stored); r != s {
					errs <- fmt.Errorf("loaded wrong value: want %d, got %#v", s, r)
				}
				continue
			}
			if r != v {
				errs <- fmt.Errorf("stored wrong value: want %d, got %#v", v, r)
			}
			atomic.AddInt32(&stores, 1)
			atomic.CompareAndSwapInt32(&stored, -1, v)
		}
		errs <- nil
	}
	for i := 0; i < n; i++ {
		go f(int32(i))
	}
	close(ch)
	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Error(err)
			i-- // don't count errors, only nils
		}
	}
}

func TestDelete(t *testing.T) {
	// Test that deleting with a nil m.v doesn't panic.
	t.Run("empty", func(t *testing.T) {
		var m generics.TopicsMap
		m.Delete("")
	})
	// Test that storing and deleting a value removes it.
	t.Run("semantic", func(t *testing.T) {
		var m generics.TopicsMap
		m.Store("k", 0)
		m.Delete("k")
		v, ok := m.Load("k")
		if ok {
			t.Errorf("deleted key k=0 was loaded, returning %v", v)
		}
		if v != nil {
			t.Errorf("deleted key k=0 returned non-nil value %v", v)
		}
	})
	// Test that concurrent reads to a value eventually see a delete.
	t.Run("concurrent", func(t *testing.T) {
		var m generics.TopicsMap
		start := make(chan bool)
		read := make(chan bool)
		errs := make(chan error, 1)
		m.Store("k", 0)
		go func() {
			<-start // receive start signal
			_, ok := m.Load("k")
			close(read) // signal that we've read
			if !ok {
				errs <- fmt.Errorf("initial load failed")
				return
			}
			for i := 0; i < 1e7; i++ {
				_, ok = m.Load("k")
				if !ok {
					errs <- nil
					return
				}
				runtime.Gosched()
			}
			errs <- fmt.Errorf("too many iterations without seeing delete")
		}()
		close(start) // send start signal
		<-read       // receive read signal
		m.Delete("k")
		if err := <-errs; err != nil {
			t.Error(err)
		}
	})
}

func TestLoadAndDelete(t *testing.T) {
	// Test that deleting with a nil m.v returns nil.
	t.Run("empty", func(t *testing.T) {
		var m generics.TopicsMap
		r, ok := m.LoadAndDelete("")
		if ok {
			t.Errorf("empty map loaded %#v", r)
		}
		if r != nil {
			t.Errorf("empty map returned %#v on delete", r)
		}
	})
	// Test that storing and deleting a value removes it.
	t.Run("semantic", func(t *testing.T) {
		var m generics.TopicsMap
		m.Store("k", 0)
		v, ok := m.LoadAndDelete("k")
		if !ok {
			t.Errorf("delete indicated no value")
		}
		if v != 0 {
			t.Errorf("delete returned wrong value: want 0/true, got %v/%t", v, ok)
		}
		v, ok = m.Load("k")
		if ok {
			t.Errorf("deleted key k=0 was loaded, returning %v", v)
		}
		if v != nil {
			t.Errorf("deleted key k=0 returned non-nil value %v", v)
		}
	})
	// Test that concurrent reads to a value eventually see a delete.
	t.Run("concurrent", func(t *testing.T) {
		var m generics.TopicsMap
		start := make(chan bool)
		read := make(chan bool)
		errs := make(chan error, 1)
		m.Store("k", 0)
		go func() {
			<-start // receive start signal
			_, ok := m.Load("k")
			close(read) // signal that we've read
			if !ok {
				errs <- fmt.Errorf("initial load failed")
				return
			}
			for i := 0; i < 1e7; i++ {
				_, ok = m.Load("k")
				if !ok {
					errs <- nil
					return
				}
				runtime.Gosched()
			}
			errs <- fmt.Errorf("too many iterations without seeing delete")
		}()
		close(start) // send start signal
		<-read       // receive read signal
		r, ok := m.LoadAndDelete("k")
		if !ok {
			t.Errorf("concurrent delete failed")
		}
		if r != 0 {
			t.Errorf("concurrent delete gave wrong result: want 0, got %v", r)
		}
		if err := <-errs; err != nil {
			t.Error(err)
		}
	})
	// Test that concurrent stores and deletes see each other.
	t.Run("write", func(t *testing.T) {
		var m generics.TopicsMap
		m.Store("k", 0)
		start := make(chan bool)
		errs := make(chan error, 1)
		n := runtime.GOMAXPROCS(0)
		f := func() {
			defer func() { errs <- nil }()
			seen := -1
			<-start
			for i := 0; i < 1e7; i++ {
				r, ok := m.LoadAndDelete("k")
				if !ok {
					continue
				}
				x := r.(int)
				if x == seen {
					m.Store("k", x)
					runtime.Gosched()
					continue
				}
				seen = x + 1
				m.Store("k", seen)
				if x >= 1e3 {
					return
				}
			}
			errs <- fmt.Errorf("too many iterations")
		}
		for i := 0; i < n; i++ {
			go f()
		}
		close(start)
		for i := 0; i < n; i++ {
			if err := <-errs; err != nil {
				t.Error(err)
				i-- // only count nils
			}
		}
	})
}

func TestRange(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		var m generics.TopicsMap
		keys := strings.Fields("a b c d e f g h i j k l m n o p q r s t u v w x y z")
		for _, k := range keys {
			m.Store(generics.EventID(k), k)
		}
		seen := make(map[generics.EventID]bool, len(keys))
		for i := 0; i < 1000; i++ {
			m.Range(func(key generics.EventID, value interface{}) bool {
				if seen[key] {
					t.Errorf("saw key %v again", key)
				}
				seen[key] = true
				return true
			})
			if len(seen) != len(keys) {
				t.Errorf("wrong keys: want %v, got %v", keys, seen)
			}
			for k := range seen {
				delete(seen, k)
			}
		}
	})
	t.Run("happens-before", func(t *testing.T) {
		var m generics.TopicsMap
		for i := 1; i <= 100; i++ {
			m.Store(generics.EventID(rune(i)), i)
			sum := 0
			m.Range(func(key generics.EventID, value interface{}) bool {
				sum += value.(int)
				return true
			})
			if sum != i*(i+1)/2 {
				t.Errorf("wrong values: sum should be %d, got %d", i*(i+1)/2, sum)
			}
		}
	})
	t.Run("stop", func(t *testing.T) {
		var m generics.TopicsMap
		keys := strings.Fields("a b c d e f g h i j k l m n o p q r s t u v w x y z")
		for _, k := range keys {
			m.Store(generics.EventID(k), k)
		}
		seen := make([]generics.EventID, 0, 1)
		m.Range(func(key generics.EventID, value interface{}) bool {
			seen = append(seen, key)
			return false
		})
		if len(seen) != 1 {
			t.Errorf("wrong number of iters: want 1, got %d (%v)", len(seen), seen)
		}
	})
	t.Run("loop-use", func(t *testing.T) {
		var m generics.TopicsMap
		keys := strings.Fields("a b c d e f g h i j k l m n o p q r s t u v w x y z")
		for _, k := range keys {
			m.Store(generics.EventID(k), k)
		}
		seen := make(map[generics.EventID]bool, len(keys))
		m.Range(func(key generics.EventID, value interface{}) bool {
			m.Load("") // should not deadlock
			m.Range(func(key generics.EventID, value interface{}) bool {
				if seen[key] {
					t.Errorf("saw key %v again", key)
				}
				seen[key] = true
				return true
			})
			if len(seen) != len(keys) {
				t.Errorf("wrong keys: want %v, got %v", keys, seen)
			}
			for k := range seen {
				delete(seen, k)
			}
			return true
		})
	})
}
