package models

import (
	"context"
	"net/http"
	"os"
)

type mongokey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// SessionFromContext returns the mongo session from the context
func SessionFromContext(c context.Context) Session {
	return c.Value(MongoKey).(Session)
}

// DataLayerFromContext takes a request argument and return the extracted *mgo.session.
func DataLayerFromContext(r *http.Request) DataLayer {
	return SessionFromContext(r.Context()).Copy().DB(os.Getenv("DB_NAME"))
}
