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

// GetFlags fetches flags from db
func (s *Store) GetFlags(t models.Tenant) ([]models.Flag, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(os.Getenv("DB_NAME"))

	var flags []models.Flag

	err := d.C("flags").Find(bson.M{"tenant": t.ID}).All(&flags)

	if err != nil {
		return flags, err
	}
	return flags, nil

}

// GetSegments fetches segments from db
func (s *Store) GetSegments(t models.Tenant) ([]models.Segment, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(os.Getenv("DB_NAME"))
	var segments []models.Segment
	err := d.C("segments").Find(bson.M{"tenant": t.ID}).All(&segments)

	if err != nil {
		logrus.Error("Cant find segments with this tenant", err)
		return segments, err
	}
	return segments, nil
}

// GetUsers fetches segments from db
func (s *Store) GetUsers(t models.Tenant) ([]models.User, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(os.Getenv("DB_NAME"))
	var users []models.User
	err := d.C("users").Find(bson.M{"tenant": t.ID}).All(&users)

	if err != nil {
		return users, err
	}
	return users, nil
}
