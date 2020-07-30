package mongo

import (
	"toggle/server/pkg/models"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// InsertFlag saves a flag to mongo
func (s *Store) InsertFlag(f *models.Flag) error {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
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

	d := sess.DB(s.DBName)
	err := d.C("segments").Insert(seg)

	if err != nil {
		logrus.Warning(err)
		return err
	}

	return nil
}

//InsertUser saves a user to mongo
func (s *Store) InsertUser(u *models.User) error {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	info, err := d.C("users").Upsert(bson.M{"key": u.Key}, u)

	if err != nil {
		logrus.Warning(err)
		return err
	}
	logrus.Println("info", info.UpsertedId)
	return nil
}

// UpsertUser retreives or registers user from db
func (s *Store) UpsertUser(u *models.User) (*models.User, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)

	_, err := d.C("users").Upsert(bson.M{"key": u.Key}, u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// InsertAttributes adds custom user attributes to db
func (s *Store) InsertAttributes(a []models.Attribute) error {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)

	for _, attr := range a {
		_, err := d.C("attributes").Upsert(bson.M{"name": attr.Name}, attr)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFlags fetches flags from db
func (s *Store) GetFlags(t models.Tenant) ([]models.Flag, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)

	var flags []models.Flag

	err := d.C("flags").Find(bson.M{"tenant": t.ID}).All(&flags)

	if err != nil {
		return flags, err
	}
	return flags, nil

}

// GetFlag retreives a single flag given a key
func (s *Store) GetFlag(key string) (*models.Flag, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	var flag models.Flag
	err := d.C("flags").Find(bson.M{"key": key}).One(&flag)
	if err != nil {
		return nil, err
	}
	return &flag, nil
}

// GetSegments fetches segments from db
func (s *Store) GetSegments(t models.Tenant) ([]models.Segment, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
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

	d := sess.DB(s.DBName)
	var users []models.User
	err := d.C("users").Find(bson.M{"tenant": t.ID}).All(&users)

	if err != nil {
		return users, err
	}
	return users, nil
}

// GetTenant finds the tenant based on key ie an email
func (s *Store) GetTenant(key string) *models.Tenant {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	var t models.Tenant

	err := d.C("tenants").Find(bson.M{"key": key}).One(&t)
	if err != nil {
		return nil
	}
	return &t
}

// InsertTenant adds a tenant to db
func (s *Store) InsertTenant(t *models.Tenant) error {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	err := d.C("tenants").Insert(t)
	if err != nil {
		logrus.Warning(err)
		return err
	}

	return nil
}
