package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/sebasromero/shortenerUrl/internal/types"
)

var CreateUrlCache = NewCache[string, types.InsertUrl]()
var GetUrlCache = NewCache[string, types.LongUrlResponse]()

type item[V any] struct {
	value  V
	expiry time.Time
}

func (i item[V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

type Cache[K comparable, V any] struct {
	items map[K]item[V]
	mu    sync.Mutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	cache := &Cache[K, V]{
		items: make(map[K]item[V]),
	}

	go func() {
		for range time.Tick(5 * time.Minute) {
			cache.mu.Lock()
			fmt.Println("Checked", cache.items)
			for key, item := range cache.items {
				if item.isExpired() {
					delete(cache.items, key)
				}
			}
			cache.mu.Unlock()
		}
	}()

	return cache
}

func (cache *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.items[key] = item[V]{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

func (cache *Cache[K, V]) Get(key K) (V, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	item, found := cache.items[key]
	if !found {
		return item.value, false
	}

	if item.isExpired() {
		delete(cache.items, key)
		return item.value, false
	}

	return item.value, true
}

func (cache *Cache[K, V]) Remove(key K) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	delete(cache.items, key)
}

func (cache *Cache[K, V]) Pop(key K) (V, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	item, found := cache.items[key]
	if !found {
		return item.value, false
	}

	delete(cache.items, key)

	if item.isExpired() {
		return item.value, false
	}

	return item.value, true
}
