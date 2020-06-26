package models

import "gopkg.in/mgo.v2/bson"

// Flag represents a feature flag object
type Flag struct {
	ID         bson.ObjectId   `bson:"_id" json:"id"`
	Name       string          `json:"name" bson:"name"`
	Key        string          `json:"key" bson:"key"`
	Enabled    bool            `bson:"enabled" json:"enabled"`
	Variations []Variation     `bson:"variations" json:"variations"`
	Users      []bson.ObjectId `bson:"users,omitempty" json:"users,omitempty"`
	Targets    []Target        `bson:"targets,omitempty" json:"targets,omitempty"`
}

// Rule is a constraint placed on users being evaluated
type Rule struct {
	Attribute string `json:"attribute" bson:"attribute"`
	Operator  string `json:"operator,omitempty" bson:"operator,omitempty"`
	Value     string `json:"value" bson:"value"`
}

// Segment represents a specific group of users
type Segment struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Rules []Rule
	Users []User
}

//Variation represents a toggle option for a flag
type Variation struct {
	Name    string `bson:"name" json:"name"`
	Percent int16  `bson:"percent" json:"percent"`
}

// User represents a client request
type User struct {
	ID         bson.ObjectId          `bson:"_id" json:"id"`
	Key        string                 `json:"key" bson:"key"`
	Attributes map[string]interface{} `json:"attributes,omitempty" bson:"attributes,omitempty"`
	Tenant     bson.ObjectId          `json:"tenant" bson:"tenant"`
}

// Target is a specific user constraint
type Target struct {
	Rule       Rule         `bson:"rule" json:"rule"`
	Variations *[]Variation `bson:"variations" json:"variations"`
}

// Tenant is a user of the system
type Tenant struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
}
