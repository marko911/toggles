package mongo

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

type MockDatabase struct{}

type MockCollection struct{}

// Find mock.
func (mc MockCollection) Find(query interface{}) *mgo.Query {
	return nil
}

// FindId mock.
func (mc MockCollection) FindId(id interface{}) *mgo.Query {
	return nil
}

// Count mock.
func (mc MockCollection) Count() (n int, err error) {
	return 10, nil
}

// Insert mock.
func (mc MockCollection) Insert(docs ...interface{}) error {
	return nil
}

// Remove mock.
func (mc MockCollection) Remove(selector interface{}) error {
	return nil
}

// Update mock.
func (mc MockCollection) Update(selector interface{}, update interface{}) error {
	return nil
}

// Upsert mock.
func (mc MockCollection) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return nil, nil
}

//Pipe mock
func (mc MockCollection) Pipe(pipeline interface{}) *mgo.Pipe {
	return nil
}

// EnsureIndex mock.
func (mc MockCollection) EnsureIndex(index mgo.Index) error {
	return nil
}

// RemoveAll mock.
func (mc MockCollection) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	return nil, nil
}

func (mc MockCollection) UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return nil, nil
}

// C mocks mgo.Database(name).Collection(name).
func (db MockDatabase) C(name string) Collection {
	return MockCollection{}
}

// MockSession satisfies Session and act as a mock of *mgo.session.
type MockSession struct{}

// NewMockSession mosck NewSession.
func NewMockSession() Session {
	return MockSession{}
}

func NewMockStore() Store {
	return Store{DBName: "mock"}
}

// Close mocks mgo.Session.Close().
func (ms MockSession) Close() {}

// Copy mocks mgo.Session.Copy().
// Regarding the context of use, no need to actually Copy the mock.
func (ms MockSession) Copy() Session {
	return ms
}

// DB mocks mgo.Session.DB().
func (ms MockSession) DB(name string) DataLayer {
	db := MockDatabase{}
	return db
}

// SetSafe mocks mgo.Session.SetSafe().
func (ms MockSession) SetSafe(safe *mgo.Safe) {}

// SetSyncTimeout mocks mgo.Session.SetSyncTimeout().
func (ms MockSession) SetSyncTimeout(d time.Duration) {}

// SetSocketTimeout mocks mgo.Session.SetSocketTimeout().
func (ms MockSession) SetSocketTimeout(d time.Duration) {}
