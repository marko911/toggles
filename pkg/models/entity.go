package models

import "gopkg.in/mgo.v2/bson"

// Flag represents a feature flag object
type Flag struct {
	ID         bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Name       string          `json:"name" bson:"name"`
	Key        string          `json:"key" bson:"key"`
	Enabled    bool            `json:"enabled" bson:"enabled"`
	Variations []Variation     `json:"variations" bson:"variations"`
	Users      []bson.ObjectId `json:"users,omitempty" bson:"users,omitempty"`
	Targets    []Target        `json:"targets,omitempty" bson:"targets,omitempty"`
	Tenant     bson.ObjectId   `json:"tenant" bson:"tenant"`
}

// Rule is a constraint placed on users being evaluated
type Rule struct {
	Attribute string `json:"attribute" bson:"attribute"`
	Operator  string `json:"operator,omitempty" bson:"operator,omitempty"`
	Value     string `json:"value" bson:"value"`
}

// Segment represents a specific group of users
type Segment struct {
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name" bson:"name"`
	Key    string        `json:"key" bson:"key"`
	Rules  []Rule        `json:"rules,omitempty" bson:"rules,omitempty"`
	Users  []string      `json:"users,omitempty" bson:"users,omitempty"` // user keys
	Tenant bson.ObjectId `json:"tenant" bson:"tenant"`
}

//Variation represents a toggle option for a flag
type Variation struct {
	Name    string `json:"name" bson:"name"`
	Percent int16  `json:"percent" bson:"percent"`
}

// User represents a client request
type User struct {
	ID         bson.ObjectId          `json:"id,omitempty" bson:"_id,omitempty"`
	Key        string                 `json:"key" bson:"key"`
	Attributes map[string]interface{} `json:"attributes,omitempty" bson:"attributes,omitempty"`
	Tenant     bson.ObjectId          `json:"tenant,omitempty" bson:"tenant"`
}

//Attribute represents a custom user attribute ie. a user group, age, gender
type Attribute struct {
	ID     bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string        `json:"name" bson:"name"`
	Tenant bson.ObjectId `json:"tenant,omitempty" bson:"tenant"`
}

// Target is a specific user constraint
type Target struct {
	Rule       Rule         `json:"rule" bson:"rule"`
	Variations *[]Variation `json:"variations" bson:"variations"`
}

// Tenant is a user of the system
type Tenant struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}
