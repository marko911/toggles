package mongo

import (
	"errors"
	"fmt"
	"strings"
	"toggle/server/pkg/models"

	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

// InsertFlag saves a flag to mongo
func (s *Store) InsertFlag(f *models.Flag) error {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	err := d.C("flags").Insert(f)

	if err != nil {
		logrus.Warning(err)

		if strings.Contains(err.Error(), "duplicate key error") {
			return errors.New("flag key already exists")
		}
		return err
	}

	return nil
}

// UpdateFlag is for updating existing flags in mongo
func (s *Store) UpdateFlag(f *models.Flag) error {
	sess := s.Copy()
	defer sess.Close()
	d := sess.DB(s.DBName)
	err := d.C("flags").Update(bson.M{"_id": f.ID}, f)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// UpdateSegment is for updating existing segments in mongo
func (s *Store) UpdateSegment(seg *models.Segment) error {
	sess := s.Copy()
	defer sess.Close()
	d := sess.DB(s.DBName)
	err := d.C("segments").Update(bson.M{"_id": seg.ID}, seg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = d.C("flags").UpdateAll(bson.M{"targets.segment._id": seg.ID}, bson.M{"$set": bson.M{
		"targets.$.segment": seg,
	}})

	if err != nil {
		logrus.Error(err)
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
		if strings.Contains(err.Error(), "duplicate key error") {
			return errors.New("segment key already exists")
		}
		return err
	}

	return nil
}

// InsertEvaluation adds an evaluation record to mongo
func (s *Store) InsertEvaluation(e *models.Evaluation) error {
	sess := s.Copy()
	defer sess.Close()
	d := sess.DB(s.DBName)
	fmt.Println("Inserting evaluation!", e.Variation)
	err := d.C("evaluations").Insert(e)
	// err := d.C("evaluations").Update(bson.M{"flag.key": e.Flag.Key}, bson.M{"$inc": bson.M{"count": 1}})
	if err != nil {
		// insertErr := d.C("evaluations").Insert(e)
		return err
	}
	return nil
}

//InsertUser saves a user to mongo
func (s *Store) InsertUser(u *models.User) error {
	sess := s.Copy()
	defer sess.Close()
	d := sess.DB(s.DBName)
	_, err := d.C("users").Upsert(bson.M{"key": u.Key}, u)
	if err != nil {
		insertErr := d.C("users").Insert(u)
		return insertErr
	}
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
			newErr := d.C("attributes").Insert(&attr)
			return newErr
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

	err := d.C("flags").Find(bson.M{"tenant": t.ID}).Sort("-evaluated").All(&flags)

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

// GetTenantFromAPIKey gets the tenant from apiKey / client key
func (s *Store) GetTenantFromAPIKey(apiKey string) *models.Tenant {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	var t models.Tenant

	err := d.C("tenants").Find(bson.M{"apiKey": apiKey}).One(&t)
	if err != nil {
		return nil
	}
	return &t
}

// InsertTenant adds a tenant to db
func (s *Store) InsertTenant(t *models.Tenant) (*models.Tenant, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	var existing models.Tenant

	if err := d.C("tenants").Find(bson.M{"key": t.Key}).One(&existing); err != nil {
		fmt.Println("inserting new teannt", t)
		err := d.C("tenants").Insert(t)
		if err != nil {
			logrus.Warning(err)
			return nil, err
		}
		err = d.C("tenants").Find(bson.M{"key": t.Key}).One(&existing)
	}

	return &existing, nil
}

// GetEvals returns all evaluations from db
func (s *Store) GetEvals() ([]models.Evaluation, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	var evals []models.Evaluation
	err := d.C("evaluations").Find(nil).All(&evals)

	return evals, err
}

// GetFlagEvals gets all evals for single flag
func (s *Store) GetFlagEvals(id bson.ObjectId, page int, pageSize int) ([]models.Evaluation, int, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	var evals []models.Evaluation
	q := d.C("evaluations").Find(bson.M{"flag._id": id}).Sort("-createdAt")
	count, _ := q.Count()
	err := q.Skip((page - 1) * pageSize).Limit(pageSize).All(&evals)
	if err != nil {
		logrus.Error("Error adding to evals slice", err)
	}
	return evals, count, err
}

//GetFlagStats returns a breakdown of flag evalutaion stats
func (s *Store) GetFlagStats(id bson.ObjectId) (*models.FlagStats, error) {
	sess := s.Copy()
	defer sess.Close()

	d := sess.DB(s.DBName)
	var stats models.FlagStats
	pipe := d.C("evaluations").Pipe([]bson.M{
		{"$match": bson.M{"flag._id": id}},

		{"$group": bson.M{"_id": bson.M{"variation": "$variation.name"},
			"count": bson.M{"$sum": 1},
		}},
	})
	err := pipe.Iter().All(&stats.Counts)

	if len(stats.Counts) == 0 {
		pipe = d.C("flags").Pipe([]bson.M{
			{"$match": bson.M{"_id": id}},
			{"$unwind": "$variations"},
			{"$project": bson.M{

				"_id":   bson.M{"variation": "$variations.name"},
				"count": "0",
			}},
		})
		err = pipe.Iter().All(&stats.Counts)
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
