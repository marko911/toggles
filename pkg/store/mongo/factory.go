package mongo

import (
	"fmt"
	"os"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// NewMongoStore returns a new Mongo Session.
func NewMongoStore(c *cli.Context) (*Store, error) {
	var sess *mgo.Session
	logrus.Println("Dialing Mongo...")
	if os.Getenv("ENV") != "prod" {
		dinfo := &mgo.DialInfo{
			Addrs:    []string{c.String("database-address")},
			Username: c.String("mongo-username"),
			Password: c.String("mongo-password"),
		}
		mgoSession, err := mgo.DialWithInfo(dinfo)
		if err != nil {
			logrus.Fatal("error dialing local mongo", err)
		}
		sess = mgoSession

	} else {
		logrus.Println("ENV=PROD")
		url := fmt.Sprintf("mongodb://%v:%v@%v/%v?ssl=true&replicaSet=atlas-672hbg-shard-0&authSource=admin", c.String("mongo-username"), c.String("mongo-password"), c.String("database-address"), c.String("database-name"))
		// url := "mongodb://backend_api_user:l30m355i@cluster0-shard-00-00.ptlpo.gcp.mongodb.net:27017,cluster0-shard-00-01.ptlpo.gcp.mongodb.net:27017,cluster0-shard-00-02.ptlpo.gcp.mongodb.net:27017/toggles?ssl=true&replicaSet=atlas-672hbg-shard-0&authSource=admin"
		d, err := mgo.ParseURL(url)
		if err != nil {
			logrus.Fatal("error parsing mongo url", err)
		}
		mgoSession, err := mgo.DialWithInfo(d)
		if err != nil {
			logrus.Fatal("error dialing mongo ", err)
		}
		sess = mgoSession
	}
	logrus.Println("Connected to mongo")
	session := Store{sess, c.String("database-name")}
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

	indexes["segments"] = mgo.Index{
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
