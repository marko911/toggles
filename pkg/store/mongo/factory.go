package mongo

import (
	"time"

	"github.com/urfave/cli"
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

} // NewMongoStore returns a new Mongo Session.
}
