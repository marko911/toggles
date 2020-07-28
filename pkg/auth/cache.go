package auth

import (
	"context"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

var singletonTenantCacheOnce sync.Once
var cache *TenantCache

type mapCache map[string]bson.ObjectId

type TenantCache struct {
	mapCacheLock sync.RWMutex
	tenants      mapCache
}

// GetTenantCache gets the TenantCache
var GetTenantCache = func() *TenantCache {
	singletonTenantCacheOnce.Do(func() {
		ec := &TenantCache{
			tenants: make(map[string]bson.ObjectId),
		}
		cache = ec
	})
	return cache
}

// GetByAuthToken returns tenant id via token lookup
func (cache *TenantCache) GetByAuthToken(token string) *bson.ObjectId {
	cache.mapCacheLock.RLock()
	defer cache.mapCacheLock.RUnlock()
	tenant, ok := cache.tenants[token]
	if ok {
		return &tenant
	}
	return nil
}

type serviceKey string

//CacheServiceKey is context key for tenantCache
const CacheServiceKey serviceKey = "tenantCache"

// TenantCacheFromContext returns the create service from context
func TenantCacheFromContext(c context.Context) *TenantCache {
	return c.Value(CacheServiceKey).(*TenantCache)
}
