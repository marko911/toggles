package models

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Insert persists the flag caller model to database
func (flag *Flag) Insert(d DataLayer) error {
	err := d.C("flags").Insert(flag)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	return nil
}

// GetFlags returns all flags.
func (db *MongoDatabase) GetFlags(t Tenant) ([]Flag, error) {
	var flags []Flag

	err := db.C("flags").Find(bson.M{"tenant": t.ID}).All(&flags)
	if err != nil {
		return flags, err
	}
	return flags, nil
}
