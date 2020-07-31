package auth

import (
	"context"
	"sync"
	"time"
	"toggle/server/pkg/read"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

var singletonCacheOnce sync.Once
var cache *Cache

type tokenToTenantIDMap map[string]bson.ObjectId
type evalToCount map[bson.ObjectId]int

// Cache stores tenants and flag evaluation counters
type Cache struct {
	mapCacheLock    sync.RWMutex
	tenants         tokenToTenantIDMap
	evals           evalToCount
	refreshInterval time.Duration
}

// GetCache gets the Cache
var GetCache = func() *Cache {
	singletonCacheOnce.Do(func() {
		ec := &Cache{
			tenants:         make(map[string]bson.ObjectId),
			evals:           make(map[bson.ObjectId]int),
			refreshInterval: time.Second * 3,
		}
		cache = ec
	})
	return cache
}

// GetByAuthToken returns tenant id via token lookup
func (cache *Cache) GetByAuthToken(token string) *bson.ObjectId {
	cache.mapCacheLock.RLock()
	defer cache.mapCacheLock.RUnlock()
	tenant, ok := cache.tenants[token]
	if ok {
		return &tenant
	}
	return nil
}

// GetEvalCount returns count of flag evaluations
func (cache *Cache) GetEvalCount(flagID bson.ObjectId) int {
	cache.mapCacheLock.RLock()
	defer cache.mapCacheLock.RUnlock()
	counts, ok := cache.evals[flagID]
	if ok {
		return counts
	}
	return 0
}

//
func (cache *Cache) StartPollingEvals(s read.Service) {
	go func() {
		for range time.Tick(cache.refreshInterval) {
			err := ec.reloadMapCache()
			if err != nil {
				logrus.WithField("err", err).Error("reload evaluation cache error")
			}
		}
	}()
}

type serviceKey string

//CacheServiceKey is context key for tenantCache
const CacheServiceKey serviceKey = "tenantCache"

// CacheFromContext returns the create service from context
func CacheFromContext(c context.Context) *Cache {
	return c.Value(CacheServiceKey).(*Cache)
}
