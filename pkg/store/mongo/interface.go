package mongo

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

// WrapCollection wraps a mgo.Collection to embed methods in models.
type WrapCollection struct {
	*mgo.Collection
}

// Collection is an interface to access to the collection struct.
type Collection interface {
	Find(query interface{}) *mgo.Query
	Count() (n int, err error)
	FindId(id interface{}) *mgo.Query
	Insert(docs ...interface{}) error
	Remove(selector interface{}) error
	Update(selector interface{}, update interface{}) error
	Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	EnsureIndex(index mgo.Index) error
	RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error)
}

// Database wraps a mgo.Database to embed methods in models.
type Database struct {
	*mgo.Database
}

// C shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (d Database) C(name string) Collection {
	return &WrapCollection{Collection: d.Database.C(name)}
}

// DataLayer is an interface to access to the database struct
type DataLayer interface {
	C(name string) Collection
}

// Session is an interface to access to the Session struct.
type Session interface {
	DB(name string) DataLayer
	SetSafe(safe *mgo.Safe)
	SetSyncTimeout(d time.Duration)
	SetSocketTimeout(d time.Duration)
	Close()
	Copy() Session
}

// Store is a Mongo session.
type Store struct {
	*mgo.Session
	DBName string
}

// DB shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (s Store) DB(name string) DataLayer {
	return &Database{Database: s.Session.DB(name)}
}

// Copy mocks mgo.Session.Copy()
func (s Store) Copy() Session {
	return Store{s.Session.Copy(), s.DBName}
}
