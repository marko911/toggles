package message

import (
	"github.com/bugsnag/bugsnag-go"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// NewNatsClient returns a nats client connection
func NewNatsClient(c *cli.Context) *nats.EncodedConn {
	// natsServers := strings.Join(c.StringSlice("nats-server-url"), ",")
	nc, err := nats.Connect("nats://nats:4222")
	conn, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		bugsnag.Notify(err)
		logrus.Fatal(err)
	}
	return conn
}
