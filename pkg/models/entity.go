package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Flag represents a feature flag object
type Flag struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name       string        `json:"name" bson:"name"`
	Key        string        `json:"key" bson:"key"`
	Enabled    bool          `json:"enabled" bson:"enabled"`
	Variations []Variation   `json:"variations" bson:"variations"`
	Targets    []Target      `json:"targets,omitempty" bson:"targets,omitempty"`
	Tenant     bson.ObjectId `json:"tenant" bson:"tenant"`
	Limit      int           `json:"limit,omitempty" bson:"limit,omitempty"`
	Evaluated  time.Time     `json:"evaluated,omitempty" bson:"evaluated,omitempty"`
}

// Rule is a constraint placed on users being evaluated
type Rule struct {
	Attribute string `json:"attribute" bson:"attribute"`
	Operator  string `json:"operator,omitempty" bson:"operator,omitempty"`
	Value     string `json:"value" bson:"value"`
}

// Segment represents a specific group of users
type Segment struct {
	ID     bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string        `json:"name,omitempty" bson:"name,omitempty"`
	Key    string        `json:"key,omitempty" bson:"key,omitempty"`
	Rules  []Rule        `json:"rules,omitempty" bson:"rules,omitempty"`
	Users  []string      `json:"users,omitempty" bson:"users,omitempty"` // user keys
	Tenant bson.ObjectId `json:"tenant,omitempty" bson:"tenant,omitempty"`
}

//Variation represents a toggle option for a flag
type Variation struct {
	Name     string   `json:"name" bson:"name"`
	Percent  float64  `json:"percent" bson:"percent"`
	UserKeys []string `json:"users,omitempty" bson:"users,omitempty"` // if a variation has specific users targeted
	Limit    int      `json:"limit,omitempty" bson:"limit,omitempty"`
}

// User represents a client request
type User struct {
	ID         bson.ObjectId          `json:"id,omitempty" bson:"_id,omitempty"`
	Key        string                 `json:"key,omitempty" bson:"key,omitempty"`
	Country    string                 `json:"country,omitempty" bson:"country,omitempty"`
	Email      string                 `json:"email,omitempty" bson:"email,omitempty"`
	Name       string                 `json:"name,omitempty" bson:"name,omitempty"`
	IP         string                 `json:"ip,omitempty" bson:"ip,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty" bson:"attributes,omitempty"`
	Tenant     bson.ObjectId          `json:"tenant,omitempty" bson:"tenant,omitempty"`
}

//Attribute represents a custom user attribute ie. a user group, age, gender
type Attribute struct {
	ID   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name string        `json:"name" bson:"name"`
	// Tenant bson.ObjectId `json:"tenant,omitempty" bson:"tenant"`
}

// Tenant is a user of the system
type Tenant struct {
	Key    string        `json:"key" bson:"key"`
	ID     bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	APIKEY string        `json:"apiKey" bson:"apiKey"`
}

// Evaluation determines what variation user is shown, can be simple true or false
// or a specific variation of flag being evaluated
type Evaluation struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Variation *Variation    `json:"variation" bson:"variation"`
	Flag      Flag          `json:"flag" bson:"flag"`
	Count     int           `json:"count,omitempty" bson:"count"`
	User      interface{}   `json:"user" bson:"user"`
	CreatedAt time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

// FlagStats represents evaluation data for flag
type FlagStats struct {
	Counts []map[string]interface{} `json:"counts" bson:"counts"`
}
