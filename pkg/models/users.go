package models

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Insert persists the user to db
func (u *User) Insert(d DataLayer) error {
	err := d.C("users").Insert(u)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	return nil
}

// GetUsers fetches all users for tenant
func (db *MongoDatabase) GetUsers(t Tenant) ([]User, error) {
	var users []User
	err := db.C("users").Find(bson.M{"tenant": t.ID}).All(&users)
	if err != nil {
		return users, err
	}
	return users, nil
}
