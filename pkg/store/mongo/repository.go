package mongo

import (
	"os"
	"time"
	"toggle/server/pkg/models"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoCollection wraps a mgo.Collection to embed methods in models.
type MongoCollection struct {
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

// MongoDatabase wraps a mgo.Database to embed methods in models.
type MongoDatabase struct {
	*mgo.Database
}

// C shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

// DataLayer is an interface to access to the database struct
// (currently MongoDatabase).
type DataLayer interface {
	C(name string) Collection
	// GetFlags(t Tenant) ([]Flag, error)
	// GetSegments(t Tenant) ([]Segment, error)
}

// Session is an interface to access to the Session struct.
type Session interface {
	DB(name string) DataLayer
	SetSafe(safe *mgo.Safe)
	SetSyncTimeout(d time.Duration)
	SetSocketTimeout(d time.Duration)
	Close()
	Copy() Store
}

// Store is a Mongo session.
type Store struct {
	*mgo.Session
}

// DB shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (s Store) DB(name string) DataLayer {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

// Copy mocks mgo.Session.Copy()
func (s Store) Copy() Store {
	return Store{s.Session.Copy()}
}

// NewMongoStore returns a new Mongo Session.
func NewMongoStore(c *cli.Context) (*Store, error) {
	mgoSession, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{c.String("database-address")},
		Username: c.String("mongo-username"),
		Password: c.String("mongo-password"),
	})
	if err != nil {
		return nil, err
	}
	session := Store{mgoSession}
	session.SetSafe(&mgo.Safe{})
	session.SetSyncTimeout(3 * time.Second)
	session.SetSocketTimeout(3 * time.Second)
	return &session, nil

}

// InsertFlag saves a flag to mongo
func (s *Store) InsertFlag(f *models.Flag) error {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(os.Getenv("DB_NAME"))
	err := d.C("flags").Insert(f)

	if err != nil {
		logrus.Warning(err)
		return err
	}

	return nil
}

//InsertSegment saves a segment to mongo
func (s *Store) InsertSegment(seg *models.Segment) error {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(os.Getenv("DB_NAME"))
	err := d.C("segments").Insert(seg)

	if err != nil {
		logrus.Warning(err)
		return err
	}

	return nil
}

//InsertUser saves a segment to mongo
func (s *Store) InsertUser(u *models.User) error {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(os.Getenv("DB_NAME"))
	err := d.C("users").Insert(u)

	if err != nil {
		logrus.Warning(err)
		return err
	}

	return nil
}

// UpsertUser retreives or registers user from db
func (s *Store) UpsertUser(u *models.User) (*models.User, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(os.Getenv("DB_NAME"))

	_, err := d.C("users").Upsert(bson.M{"key": u.Key}, u)
	if err != nil {
		return u, err
	}
	return u, nil
}
