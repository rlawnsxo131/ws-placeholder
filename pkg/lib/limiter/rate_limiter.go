package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	*rate.Limiter

	lastSeen time.Time
}

type memoryStore struct {
	visitors map[string]*visitor
	mutex    sync.Mutex
	rate     rate.Limit // for more info check out Limiter docs - https://pkg.go.dev/golang.org/x/time/rate#Limit.

	burst       int
	expiresIn   time.Duration
	lastCleanup time.Time

	timeNow func() time.Time
}

func NewRateLimiterMemoryStore(r rate.Limit) *memoryStore {
	store := &memoryStore{}
	store.rate = r
	store.burst = int(r)
	store.expiresIn = 3 * time.Minute
	store.visitors = make(map[string]*visitor)
	store.timeNow = time.Now
	store.lastCleanup = store.timeNow()

	return store
}

func (store *memoryStore) Allow(identifier string) bool {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	limiter, exists := store.visitors[identifier]
	if !exists {
		limiter = &visitor{
			Limiter: rate.NewLimiter(store.rate, store.burst),
		}
		store.visitors[identifier] = limiter
	}

	now := store.timeNow()
	limiter.lastSeen = now
	if now.Sub(store.lastCleanup) > store.expiresIn {
		// cleanupStaleVisitors helps manage the size of the visitors map by removing stale records
		// of users who haven't visited again after the configured expiry time has elapsed
		for id, visitor := range store.visitors {
			if store.timeNow().Sub(visitor.lastSeen) > store.expiresIn {
				delete(store.visitors, id)
			}
		}
		store.lastCleanup = store.timeNow()
	}

	return limiter.AllowN(store.timeNow(), 1)
}
