package auth

import (
	"context"
	"toggle/server/pkg/models"
)

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

type mongokey string
type tenantkey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// TenantKey is context key for tenant
const TenantKey tenantkey = "tenant"

// TenantFromContext returns the tenant from context
func TenantFromContext(c context.Context) *models.Tenant {
	if val, ok := c.Value(TenantKey).(*models.Tenant); ok {
		return val
	}
	return nil
}
