package cache

import (
	"sync"
	"time"
)

// InstanceCache structure for caching instance IDs.
// Stores mapping between user IDs and instance IDs with expiration.
type InstanceCache struct {
	// key: user ID (string)
	// value: cache item (contains instance ID and expiration time)
	cache *sync.Map
}

// cacheItem represents a cached entry.
// Wraps instance ID and expiration time.
type cacheItem struct {
	instanceID string    // The instance ID to be cached
	expireAt   time.Time // Expiration time
	// Prevents data inconsistency from long-term caching
}

// NewInstanceCache initializes a new instance ID cache.
func NewInstanceCache() *InstanceCache {
	cache := &InstanceCache{
		cache: new(sync.Map),
	}

	// Start background goroutine for periodic expiration cleanup
	go func() {
		// Clean up expired cache every 5 minutes
		// Adjust frequency based on actual needs
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			cache.cleanupExpired()
		}
	}()

	return cache
}

// Set caches the instance ID for a user.
// expireSeconds: expiration time in seconds (e.g. 3600 = 1 hour)
func (c *InstanceCache) Set(userID, instanceID string, expireSeconds int64) {
	item := cacheItem{
		instanceID: instanceID,
		expireAt:   time.Now().Add(time.Duration(expireSeconds) * time.Second),
	}
	c.cache.Store(userID, item)
}

// Get retrieves the cached instance ID for a user.
// Returns: instance ID, whether it exists and is not expired
func (c *InstanceCache) Get(userID string) (string, bool) {
	val, ok := c.cache.Load(userID)
	if !ok {
		return "", false // Cache does not exist
	}

	item, ok := val.(cacheItem)
	if !ok {
		// Invalid data format, clean up dirty data
		c.cache.Delete(userID)
		return "", false
	}

	// Check if expired
	if time.Now().After(item.expireAt) {
		// Automatically clean up expired entries
		c.cache.Delete(userID)
		return "", false
	}

	return item.instanceID, true
}

// Delete actively removes the cache for a specific user.
// Useful when user's instance ID is updated.
func (c *InstanceCache) Delete(userID string) {
	c.cache.Delete(userID)
}

// cleanupExpired removes all expired cache entries.
func (c *InstanceCache) cleanupExpired() {
	c.cache.Range(func(key, value interface{}) bool {
		item, ok := value.(cacheItem)
		if !ok {
			// Invalid data format, clean up dirty data
			c.cache.Delete(key)
			return true
		}

		// Check if expired
		if time.Now().After(item.expireAt) {
			c.cache.Delete(key)
		}
		return true
	})
}
