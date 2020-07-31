package auth

import "context"

// FromContext returns Authorizer from context
func FromContext(c context.Context) *Authorizer {
	return c.Value(ServiceKey).(*Authorizer)
}

type serviceKey string

//CacheServiceKey is context key for tenantCache
const CacheServiceKey serviceKey = "tenantCache"

// CacheFromContext returns the create service from context
func CacheFromContext(c context.Context) *Cache {
	return c.Value(CacheServiceKey).(*Cache)
}
