package models

import (
	"context"
	"net/http"
	"os"
)

type mongokey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// TenantKey is context key for tenant
const TenantKey string = "tenant"

// SessionFromContext returns the mongo session from the context
func SessionFromContext(c context.Context) Session {
	return c.Value(MongoKey).(Session)
}

// TenantFromContext returns the tenant from context
func TenantFromContext(c context.Context) Tenant {
	return c.Value(TenantKey).(Tenant)
}

// DataLayerFromContext takes a request argument and return the extracted *mgo.session.
func DataLayerFromContext(r *http.Request) DataLayer {
	return SessionFromContext(r.Context()).Copy().DB(os.Getenv("DB_NAME"))
}
