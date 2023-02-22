package generics

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type CacheObject[T any] struct {
	Expiration int64
	Value      T
}

func (item *CacheObject[T]) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

const (
	NeverExpire   time.Duration = -1
	DefaultExpire time.Duration = 0
)

type janitor[T any] struct {
	Interval time.Duration
	stop     chan struct{}
}

func (j *janitor[T]) Run(c *cache[T]) {
	timer := time.NewTicker(j.Interval)
	for {
		select {
		case <-timer.C:
			c.CleanExpired()
		case <-j.stop:
			timer.Stop()
			return
		}
	}
}

func stopJanitor[T any](c *Cache[T]) {
	c.janitor.stop <- struct{}{}
}

func runJanitor[T any](cache *cache[T], interval time.Duration) {
	service := &janitor[T]{
		Interval: interval,
		stop:     make(chan struct{}),
	}
	cache.janitor = service
	go service.Run(cache)
}

func newCache[T any](deadline time.Duration, objectsMap map[string]CacheObject[T]) *cache[T] {
	if deadline == 0 {
		deadline = -1
	}
	result := &cache[T]{
		defaultExpiration: deadline,
		objects:           objectsMap,
	}
	return result
}

func newCacheWithJanitor[T any](deadline time.Duration, cleanup time.Duration, objectsMap map[string]CacheObject[T]) *Cache[T] {
	cacher := newCache[T](deadline, objectsMap)
	result := &Cache[T]{cacher}
	if cleanup > 0 {
		runJanitor(cacher, cleanup)
		runtime.SetFinalizer(result, stopJanitor[T]) // I know, I know - you can stop the garbage collection and finalizers "are for nothing"...
	}
	return result
}

func NewCache[T any](expiration, cleanup time.Duration) *Cache[T] {
	return newCacheWithJanitor[T](expiration, cleanup, make(map[string]CacheObject[T]))
}

func NewCacheFrom[T any](defaultExpiration, cleanupInterval time.Duration, objectsMap map[string]CacheObject[T]) *Cache[T] {
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, objectsMap)
}

type cache[T any] struct {
	mu                sync.RWMutex
	defaultExpiration time.Duration
	janitor           *janitor[T]
	beforeDeletion    func(string, T)
	objects           map[string]CacheObject[T]
}

type Cache[T any] struct {
	*cache[T]
}

// Set an object to the cache, replacing any existing object
func (c *cache[T]) Set(key string, object T, defaultExpiration time.Duration) {
	expiration := int64(0)
	if defaultExpiration == DefaultExpire {
		defaultExpiration = c.defaultExpiration
	}
	if defaultExpiration > 0 {
		expiration = time.Now().Add(defaultExpiration).UnixNano()
	}
	c.mu.Lock()
	c.objects[key] = CacheObject[T]{
		Value:      object,
		Expiration: expiration,
	}
	c.mu.Unlock()
}

func (c *cache[T]) set(key string, object T, defaultExpiration time.Duration) {
	expiration := int64(0)
	if defaultExpiration == DefaultExpire {
		defaultExpiration = c.defaultExpiration
	}
	if defaultExpiration > 0 {
		expiration = time.Now().Add(defaultExpiration).UnixNano()
	}
	c.objects[key] = CacheObject[T]{
		Value:      object,
		Expiration: expiration,
	}
}

// SetWithDefault adds an object to the cache, replacing any existing item, using the default expiration value.
func (c *cache[T]) SetWithDefault(key string, object T) {
	c.Set(key, object, DefaultExpire)
}

// Store an object to the cache only if an item doesn't already exist for the given key, or if the existing item has expired. Returns an error otherwise.
func (c *cache[T]) Store(key string, object T, d time.Duration) error {
	c.mu.Lock()
	_, has := c.get(key)
	if has {
		c.mu.Unlock()
		return fmt.Errorf("object %s already exists", key)
	}
	c.set(key, object, d)
	c.mu.Unlock()
	return nil
}

// Replace a new value for the cache key only if it already exists, and the existing object hasn't expired. Returns an error otherwise.
func (c *cache[T]) Replace(key string, object T, d time.Duration) error {
	c.mu.Lock()
	_, has := c.get(key)
	if !has {
		c.mu.Unlock()
		return fmt.Errorf("object %s not found", key)
	}
	c.set(key, object, d)
	c.mu.Unlock()
	return nil
}

// Read an object from the cache. Returns the item or nil, and a bool indicating whether the key was found.
func (c *cache[T]) Read(key string) (T, bool) {
	c.mu.RLock()
	object, has := c.objects[key]
	if !has {
		c.mu.RUnlock()
		return *new(T), false
	}

	if object.Expiration > 0 {
		if time.Now().UnixNano() > object.Expiration {
			c.mu.RUnlock()
			return *new(T), false
		}
	}
	c.mu.RUnlock()
	return object.Value, true
}

// GetWithExpiration returns an object and its expiration time from the cache.
// It returns the object or nil, the expiration time if one is set (if the object never expires a zero value for time.Time is returned), and a bool indicating whether the key was found.
func (c *cache[T]) GetWithExpiration(key string) (T, time.Time, bool) {
	c.mu.RLock()
	object, found := c.objects[key]
	if !found {
		c.mu.RUnlock()
		return *new(T), time.Time{}, false
	}

	if object.Expiration > 0 {
		if time.Now().UnixNano() > object.Expiration {
			c.mu.RUnlock()
			return *new(T), time.Time{}, false
		}

		c.mu.RUnlock()
		return object.Value, time.Unix(0, object.Expiration), true
	}

	c.mu.RUnlock()
	return object.Value, time.Time{}, true
}

func (c *cache[T]) get(k string) (T, bool) {
	object, has := c.objects[k]
	if !has {
		return *new(T), false
	}

	if object.Expiration > 0 {
		if time.Now().UnixNano() > object.Expiration {
			return *new(T), false
		}
	}
	return object.Value, true
}

// Delete an object from the cache. Does nothing if the key is not in the cache.
func (c *cache[T]) Delete(key string) {
	c.mu.Lock()
	v, deleted := c.delete(key)
	c.mu.Unlock()
	if deleted {
		c.beforeDeletion(key, v)
	}
}

func (c *cache[T]) delete(key string) (T, bool) {
	if c.beforeDeletion != nil {
		if object, has := c.objects[key]; has {
			delete(c.objects, key)
			return object.Value, true
		}
	}
	delete(c.objects, key)
	return *new(T), false
}

type cachePair[T any] struct {
	key   string
	value T
}

// CleanExpired all expired objects from the cache.
func (c *cache[T]) CleanExpired() {
	var evictedItems []cachePair[T]
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.objects {
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, cachePair[T]{k, ov})
			}
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		c.beforeDeletion(v.key, v.value)
	}
}

// OnBeforeDeletion registers a function to be called before eviction
func (c *cache[T]) OnBeforeDeletion(hook func(string, T)) {
	c.mu.Lock()
	c.beforeDeletion = hook
	c.mu.Unlock()
}

// Save everything in cache to an io.Writer, in gob encoding format
func (c *cache[T]) Save(w io.Writer) error {
	enc := gob.NewEncoder(w)
	c.mu.RLock()
	defer c.mu.RUnlock()

	var t T
	switch reflect.TypeOf(t).Kind() {
	case reflect.Func:
		return fmt.Errorf("error : can't encode functions")
	case reflect.Chan:
		return fmt.Errorf("error : can't encode channels")
	}

	for _, object := range c.objects {
		gob.Register(object.Value)
	}

	return enc.Encode(&c.objects)
}

// SaveFile saves cache content to a file
func (c *cache[T]) SaveFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	err = c.Save(file)
	if err != nil {
		file.Close()
		return err
	}

	return file.Close()
}

// Load fills the cache from io.Reader
func (c *cache[T]) Load(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	objects := map[string]CacheObject[T]{}
	err := decoder.Decode(&objects)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range objects {
		object, has := c.objects[key]
		if !has || object.Expired() {
			c.objects[key] = value
		}
	}
	return nil
}

// LoadFile fills the cache from a file
func (c *cache[T]) LoadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	err = c.Load(file)
	if err != nil {
		file.Close()
		return err
	}

	return file.Close()
}

// Objects lists all objects in the cache
func (c *cache[T]) Objects() map[string]CacheObject[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	objects := make(map[string]CacheObject[T], len(c.objects))
	now := time.Now().UnixNano()
	for k, v := range c.objects {
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		objects[k] = v
	}
	return objects
}

// ObjectsCount returns the number of cached objects
func (c *cache[T]) ObjectsCount() int {
	c.mu.RLock()
	n := len(c.objects)
	c.mu.RUnlock()
	return n
}

// Flush flushes the cache
func (c *cache[T]) Flush() {
	c.mu.Lock()
	c.objects = map[string]CacheObject[T]{}
	c.mu.Unlock()
}
