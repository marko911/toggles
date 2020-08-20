package mongo

import (
	"time"

	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
	mgo "gopkg.in/mgo.v2"
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
	session := Store{mgoSession, c.String("database-name")}
	session.SetSafe(&mgo.Safe{})
	session.SetSyncTimeout(3 * time.Second)
	session.SetSocketTimeout(3 * time.Second)
	return &session, nil

}

// PrepareDB sets up all mongo indexes
func PrepareDB(session Session) {
	indexes := make(map[string]mgo.Index)
	indexes["flags"] = mgo.Index{
		Key:    []string{"tenant", "key"},
		Unique: true,
	}

	for collectionName, index := range indexes {
		err := session.DB("toggles").C(collectionName).EnsureIndex(index)
		if err != nil {
			panic("Cannot ensure index.")
		}
	}
	log.Info("Prepared database indexes.")
}
