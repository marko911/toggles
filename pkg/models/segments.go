package models

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Insert persists a Segment to db
func (s *Segment) Insert(d DataLayer) error {
	err := d.C("segments").Insert(s)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	return nil
}

// GetSegments returns all segments
func (db *MongoDatabase) GetSegments(t Tenant) ([]Segment, error) {
	var segments []Segment

	err := db.C("segments").Find(bson.M{"tenant": t.ID}).All(&segments)
	if err != nil {
		return segments, err
	}
	return segments, nil
}
