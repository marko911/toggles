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

// Store stores all entities data in mongo
type Store struct {
	*mgo.Session
}

// NewStore returns a new mongo session
func NewStore(c *cli.Context) (*Store, error) {
	mgoSession, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{c.String("database-address")},
		Username: c.String("mongo-username"),
		Password: c.String("mongo-password"),
	})
	if err != nil {
		return nil, err
	}
	session := &Store{mgoSession}
	session.SetSafe(&mgo.Safe{})
	session.SetSyncTimeout(3 * time.Second)
	session.SetSocketTimeout(3 * time.Second)
	return session, nil
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
