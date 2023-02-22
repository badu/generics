package generics

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Tree struct {
	Num      int
	Children []*Tree
}

func NewTestStruct(num int) Tree {
	return Tree{Num: num}
}

func TestCache(t *testing.T) {
	myCache := NewCache[Tree](DefaultExpire, 0)

	a, has := myCache.Read("a")
	if has {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	b, has := myCache.Read("b")
	if has {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	c, has := myCache.Read("c")
	if has {
		t.Error("Getting C found value that shouldn't exist:", c)
	}

	myCache.Set("a", Tree{Num: 1}, DefaultExpire)
	myCache.Set("b", Tree{Num: 2}, DefaultExpire)

	x, has := myCache.Read("a")
	if !has {
		t.Error("a was not found while getting a2")
	}
	assert.Equal(t, 1, x.Num)

	x, has = myCache.Read("b")
	assert.True(t, has, "b was not found while getting b2")
	assert.Equal(t, 2, x.Num)
}

func TestCacheTimes(t *testing.T) {
	var has bool

	tc := NewCache[Tree](50*time.Millisecond, 1*time.Millisecond)
	tc.Set("a", NewTestStruct(1), DefaultExpire)
	tc.Set("b", NewTestStruct(2), NeverExpire)
	tc.Set("c", NewTestStruct(3), 20*time.Millisecond)
	tc.Set("d", NewTestStruct(4), 70*time.Millisecond)

	<-time.After(25 * time.Millisecond)
	_, has = tc.Read("c")
	if has {
		t.Error("Found c when it should have been automatically deleted")
	}

	<-time.After(30 * time.Millisecond)
	_, has = tc.Read("a")
	if has {
		t.Error("Found a when it should have been automatically deleted")
	}

	_, has = tc.Read("b")
	if !has {
		t.Error("Did not find b even though it was set to never expire")
	}

	_, has = tc.Read("d")
	if !has {
		t.Error("Did not find d even though it was set to expire later than the default")
	}

	<-time.After(20 * time.Millisecond)
	_, has = tc.Read("d")
	if has {
		t.Error("Found d when it should have been automatically deleted (later than the default)")
	}
}

func TestNewCacheFrom(t *testing.T) {

	m := map[string]CacheObject[int]{
		"a": CacheObject[int]{
			Value:      1,
			Expiration: 0,
		},
		"b": CacheObject[int]{
			Value:      2,
			Expiration: 0,
		},
	}
	tc := NewCacheFrom(DefaultExpire, 0, m)
	a, has := tc.Read("a")
	if !has {
		t.Fatal("Did not find a")
	}
	if a != 1 {
		t.Fatal("a is not 1")
	}
	b, has := tc.Read("b")
	if !has {
		t.Fatal("Did not find b")
	}
	if b != 2 {
		t.Fatal("b is not 2")
	}
}

func TestStorePointerToStruct(t *testing.T) {
	tc := NewCache[*Tree](DefaultExpire, 0)
	tc.Set("foo", &Tree{Num: 1}, DefaultExpire)
	x, has := tc.Read("foo")
	if !has {
		t.Fatal("*Tree was not found for foo")
	}
	foo := x
	foo.Num++

	y, has := tc.Read("foo")
	if !has {
		t.Fatal("*Tree was not found for foo (second time)")
	}
	bar := y
	if bar.Num != 2 {
		t.Fatal("Tree.Num is not 2")
	}
}

func TestStoreInCache(t *testing.T) {
	tc := NewCache[string](DefaultExpire, 0)
	err := tc.Store("foo", "bar", DefaultExpire)
	if err != nil {
		t.Error("Couldn't add foo even though it shouldn't exist")
	}
	err = tc.Store("foo", "baz", DefaultExpire)
	if err == nil {
		t.Error("Successfully added another foo when it should have returned an error")
	}
}

func TestReplaceInCache(t *testing.T) {
	tc := NewCache[string](DefaultExpire, 0)
	err := tc.Replace("foo", "bar", DefaultExpire)
	if err == nil {
		t.Error("Replaced foo when it shouldn't exist")
	}
	tc.Set("foo", "bar", DefaultExpire)
	err = tc.Replace("foo", "bar", DefaultExpire)
	if err != nil {
		t.Error("Couldn't replace existing key foo")
	}
}

func TestDeleteFromCache(t *testing.T) {
	tc := NewCache[string](DefaultExpire, 0)
	tc.Set("foo", "bar", DefaultExpire)
	tc.Delete("foo")
	_, found := tc.Read("foo")
	if found {
		t.Error("foo was found, but it should have been deleted")
	}
}

func TestCacheItemCount(t *testing.T) {
	tc := NewCache[string](DefaultExpire, 0)
	tc.Set("foo", "1", DefaultExpire)
	tc.Set("bar", "2", DefaultExpire)
	tc.Set("baz", "3", DefaultExpire)
	if n := tc.ObjectsCount(); n != 3 {
		t.Errorf("Item count is not 3: %d", n)
	}
}

func TestFlushCache(t *testing.T) {
	tc := NewCache[string](DefaultExpire, 0)
	tc.Set("foo", "bar", DefaultExpire)
	tc.Set("baz", "yes", DefaultExpire)
	tc.Flush()
	_, found := tc.Read("foo")
	if found {
		t.Error("foo was found, but it should have been deleted")
	}

	_, found = tc.Read("baz")
	if found {
		t.Error("baz was found, but it should have been deleted")
	}
}

func TestOnCacheEvicted(t *testing.T) {
	tc := NewCache[int](DefaultExpire, 0)
	tc.Set("foo", 3, DefaultExpire)
	if tc.beforeDeletion != nil {
		t.Fatal("tc.beforeDeletion is not nil")
	}
	works := false
	tc.OnBeforeDeletion(func(k string, v int) {
		if k == "foo" && v == 3 {
			works = true
		}
		tc.Set("bar", 4, DefaultExpire)
	})
	tc.Delete("foo")
	x, _ := tc.Read("bar")
	if !works {
		t.Error("works bool not true")
	}
	if x != 4 {
		t.Error("bar was not 4")
	}
}

func TestCacheSerialization(t *testing.T) {
	tc := NewCache[Tree](DefaultExpire, 0)
	testFillAndSerialize(t, tc)
	testFillAndSerialize(t, tc)
}

func testFillAndSerialize(t *testing.T, tc *Cache[Tree]) {
	tc.Set("*struct", Tree{Num: 1}, DefaultExpire)
	tc.Set("structception", Tree{
		Num: 42,
		Children: []*Tree{
			&Tree{Num: 6174},
			&Tree{Num: 4716},
		},
	}, DefaultExpire)

	tc.Set("structceptionexpire", Tree{
		Num: 42,
		Children: []*Tree{
			&Tree{Num: 6174},
			&Tree{Num: 4716},
		},
	}, 1*time.Millisecond)

	b := &bytes.Buffer{}
	err := tc.Save(b)
	if err != nil {
		t.Fatal("couldn't save cache to buffer:", err)
	}

	oc := NewCache[Tree](DefaultExpire, 0)
	err = oc.Load(b)
	if err != nil {
		t.Fatal("couldn't load cache from buffer:", err)
	}

	<-time.After(5 * time.Millisecond)
	_, found := oc.Read("structceptionexpire")
	if found {
		t.Error("expired was found")
	}

	s1, found := oc.Read("*struct")
	if !found {
		t.Error("*struct was not found")
	}
	if s1.Num != 1 {
		t.Error("*struct.Num is not 1")
	}

	s4, found := oc.get("structception")
	if !found {
		t.Error("structception was not found")
	}
	s4r := s4
	if len(s4r.Children) != 2 {
		t.Error("Length of s4r.Children is not 2")
	}
	if s4r.Children[0].Num != 6174 {
		t.Error("s4r.Children[0].Num is not 6174")
	}
	if s4r.Children[1].Num != 4716 {
		t.Error("s4r.Children[1].Num is not 4716")
	}
}

func TestCacheInFileSerialization(t *testing.T) {
	tc := NewCache[string](DefaultExpire, 0)
	tc.Store("a", "a", DefaultExpire)
	tc.Store("b", "b", DefaultExpire)
	f, err := ioutil.TempFile("", "go-cache-cache.dat")
	if err != nil {
		t.Fatal("Couldn't create cache file:", err)
	}
	fname := f.Name()
	f.Close()
	tc.SaveFile(fname)

	oc := NewCache[string](DefaultExpire, 0)
	oc.Store("a", "aa", 0) // this should not be overwritten
	err = oc.LoadFile(fname)
	if err != nil {
		t.Error(err)
	}
	a, found := oc.Read("a")
	if !found {
		t.Error("a was not found")
	}
	astr := a
	if astr != "aa" {
		if astr == "a" {
			t.Error("a was overwritten")
		} else {
			t.Error("a is not aa")
		}
	}
	b, found := oc.Read("b")
	if !found {
		t.Error("b was not found")
	}
	if b != "b" {
		t.Error("b is not b")
	}
}

func TestCacheSerializeUnserializable(t *testing.T) {
	tc := NewCache[chan bool](DefaultExpire, 0)
	ch := make(chan bool, 1)
	ch <- true
	tc.Set("chan", ch, DefaultExpire)
	fp := &bytes.Buffer{}
	err := tc.Save(fp) // this should fail gracefully
	if assert.Error(t, err) {
		assert.NotEqual(t, err.Error(), "gob NewTypeObject can't handle type: chan bool", "Error from Save was not gob NewTypeObject can't handle type chan bool:", err)
	}
}

func BenchmarkCacheGetExpiring(b *testing.B) {
	benchmarkCacheGet(b, 5*time.Minute)
}

func BenchmarkCacheGetNotExpiring(b *testing.B) {
	benchmarkCacheGet(b, NeverExpire)
}

func benchmarkCacheGet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := NewCache[string](exp, 0)
	tc.Set("foo", "bar", DefaultExpire)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Read("foo")
	}
}

func BenchmarkRWMutexMapGet(b *testing.B) {
	b.StopTimer()
	m := map[string]string{
		"foo": "bar",
	}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m["foo"]
		mu.RUnlock()
	}
}

func BenchmarkRWMutexInterfaceMapGetStruct(b *testing.B) {
	b.StopTimer()
	s := struct{ name string }{name: "foo"}
	m := map[interface{}]string{
		s: "bar",
	}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m[s]
		mu.RUnlock()
	}
}

func BenchmarkRWMutexInterfaceMapGetString(b *testing.B) {
	b.StopTimer()
	m := map[interface{}]string{
		"foo": "bar",
	}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m["foo"]
		mu.RUnlock()
	}
}

func BenchmarkCacheGetConcurrentExpiring(b *testing.B) {
	benchmarkCacheGetConcurrent(b, 5*time.Minute)
}

func BenchmarkCacheGetConcurrentNotExpiring(b *testing.B) {
	benchmarkCacheGetConcurrent(b, NeverExpire)
}

func benchmarkCacheGetConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := NewCache[string](exp, 0)
	tc.Set("foo", "bar", DefaultExpire)
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				tc.Read("foo")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkRWMutexMapGetConcurrent(b *testing.B) {
	b.StopTimer()
	m := map[string]string{
		"foo": "bar",
	}
	mu := sync.RWMutex{}
	wg := new(sync.WaitGroup)
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < each; j++ {
				mu.RLock()
				_ = m["foo"]
				mu.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCacheGetManyConcurrentExpiring(b *testing.B) {
	benchmarkCacheGetManyConcurrent(b, 5*time.Minute)
}

func BenchmarkCacheGetManyConcurrentNotExpiring(b *testing.B) {
	benchmarkCacheGetManyConcurrent(b, NeverExpire)
}

func benchmarkCacheGetManyConcurrent(b *testing.B, exp time.Duration) {
	b.StopTimer()
	n := 10000
	tc := NewCache[string](exp, 0)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		k := "foo" + strconv.Itoa(i)
		keys[i] = k
		tc.Set(k, "bar", DefaultExpire)
	}
	each := b.N / n
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for _, v := range keys {
		go func(k string) {
			for j := 0; j < each; j++ {
				tc.Read(k)
			}
			wg.Done()
		}(v)
	}
	b.StartTimer()
	wg.Wait()
}

func BenchmarkCacheSetExpiring(b *testing.B) {
	benchmarkCacheSet(b, 5*time.Minute)
}

func BenchmarkCacheSetNotExpiring(b *testing.B) {
	benchmarkCacheSet(b, NeverExpire)
}

func benchmarkCacheSet(b *testing.B, exp time.Duration) {
	b.StopTimer()
	tc := NewCache[string](exp, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set("foo", "bar", DefaultExpire)
	}
}

func BenchmarkRWMutexMapSet(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m["foo"] = "bar"
		mu.Unlock()
	}
}

func BenchmarkCacheSetDelete(b *testing.B) {
	b.StopTimer()
	tc := NewCache[string](DefaultExpire, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set("foo", "bar", DefaultExpire)
		tc.Delete("foo")
	}
}

func BenchmarkRWMutexMapSetDelete(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m["foo"] = "bar"
		mu.Unlock()
		mu.Lock()
		delete(m, "foo")
		mu.Unlock()
	}
}

func BenchmarkCacheSetDeleteSingleLock(b *testing.B) {
	b.StopTimer()
	tc := NewCache[string](DefaultExpire, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.mu.Lock()
		tc.set("foo", "bar", DefaultExpire)
		tc.delete("foo")
		tc.mu.Unlock()
	}
}

func BenchmarkRWMutexMapSetDeleteSingleLock(b *testing.B) {
	b.StopTimer()
	m := map[string]string{}
	mu := sync.RWMutex{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m["foo"] = "bar"
		delete(m, "foo")
		mu.Unlock()
	}
}

func BenchmarkDeleteExpiredLoop(b *testing.B) {
	b.StopTimer()
	tc := NewCache[string](5*time.Minute, 0)
	tc.mu.Lock()
	for i := 0; i < 100000; i++ {
		tc.set(strconv.Itoa(i), "bar", DefaultExpire)
	}
	tc.mu.Unlock()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.CleanExpired()
	}
}

func TestGetWithExpiration(t *testing.T) {
	tc := NewCache[int](DefaultExpire, 0)

	a, expiration, found := tc.GetWithExpiration("a")
	if found || !expiration.IsZero() {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	b, expiration, found := tc.GetWithExpiration("b")
	if found || !expiration.IsZero() {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	c, expiration, found := tc.GetWithExpiration("c")
	if found || !expiration.IsZero() {
		t.Error("Getting C found value that shouldn't exist:", c)
	}

	tc.Set("a", 1, DefaultExpire)
	tc.Set("b", 2, DefaultExpire)
	tc.Set("c", 3, DefaultExpire)
	tc.Set("d", 4, NeverExpire)
	tc.Set("e", 5, 50*time.Millisecond)

	x, expiration, found := tc.GetWithExpiration("a")
	assert.True(t, found, "Didn't find a")
	assert.Equal(t, 1, x)
	assert.True(t, expiration.IsZero(), "expiration for a is not a zeroed time")

	x, expiration, found = tc.GetWithExpiration("b")
	assert.True(t, found, "Didn't find b")
	assert.Equal(t, 2, x)
	assert.True(t, expiration.IsZero(), "expiration for b is not a zeroed time")

	x, expiration, found = tc.GetWithExpiration("c")
	assert.True(t, found, "Didn't find c")
	assert.Equal(t, 3, x)
	assert.True(t, expiration.IsZero(), "expiration for c is not a zeroed time")

	x, expiration, found = tc.GetWithExpiration("d")
	assert.True(t, found, "Didn't find d")
	assert.Equal(t, 4, x)
	assert.True(t, expiration.IsZero(), "expiration for d is not a zeroed time")

	x, expiration, found = tc.GetWithExpiration("e")
	assert.True(t, found, "Didn't find e")
	assert.Equal(t, 5, x)
	assert.Equal(t, expiration.UnixNano(), tc.objects["e"].Expiration, "expiration for e is not the correct time")
	assert.Greater(t, expiration.UnixNano(), time.Now().UnixNano(), "expiration for e is in the past")
}
