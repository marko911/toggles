package models

import (
	"context"
)

type mongokey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// TenantKey is context key for tenant
const TenantKey string = "tenant"

// TenantFromContext returns the tenant from context
func TenantFromContext(c context.Context) Tenant {
	return c.Value(TenantKey).(Tenant)
}
