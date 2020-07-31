package message

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// NewNatsClient returns a nats client connection
func NewNatsClient(c *cli.Context) *nats.EncodedConn {
	natsServers := strings.Join(c.StringSlice("nats-server-url"), ",")
	nc, err := nats.Connect(natsServers)
	conn, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		logrus.Fatal(err)
	}
	return conn
}
