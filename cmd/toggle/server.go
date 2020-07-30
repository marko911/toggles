package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/create"
	"toggle/server/pkg/evaluate"
	"toggle/server/pkg/handler"
	"toggle/server/pkg/message"
	"toggle/server/pkg/read"
	"toggle/server/pkg/store/mongo"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Server for toggle API
type Server struct {
	router handler.Router
	store  *mongo.Store
}

var fatalErr error

func fatal(e error) {
	logrus.Fatal(e)
	flag.PrintDefaults()
	fatalErr = e
}

// StoreType defines available storage types
type StoreType int

const (
	// Mongo will use mongo as the store method
	Mongo StoreType = iota
)

// NewServer builds an http.Server with a router and port
func NewServer(c *cli.Context) *Server {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	s, err := mongo.NewMongoStore(c)

	if err != nil {
		logrus.Fatal(err)
		return nil
	}

	mongo.PrepareDB(s)
	create := create.NewService(s)
	read := read.NewService(s)
	evaluate := evaluate.NewService(s)

	natsClient := message.NewNatsClient(c)
	messenger := message.NewNatsService(natsClient)

	r := handler.Router{Create: create, Read: read, Message: messenger, Evaluate: evaluate, Authorizer: &auth.Authorizer{}, TenantCache: auth.GetTenantCache()}

	return &Server{r, s}
}

//Start the server
func (s Server) Start(c *cli.Context) {
	defer s.store.Close()

	logrus.SetLevel(logrus.InfoLevel)
	addr := c.String("server-address")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	natsServers := strings.Join(c.StringSlice("nats-server-url"), ",")
	nc, err := nats.Connect(natsServers)
	if err != nil {
		logrus.Fatal(err)
	}

	// Simple Publisher
	err = nc.Publish("hello", []byte("Hello asd"))
	if err != nil {
		logrus.Fatal("Failed publuish", err)
	}
	nc.Subscribe("request", func(m *nats.Msg) {
		m.Respond([]byte("answer is 42"))
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router.Handler(c),
	}
	go func() {
		logrus.Fatal(srv.ListenAndServe())
	}()

	logrus.Printf("Server Started at to http://localhost%s/", addr)

	// Wait for the interrupt signal.
	<-interrupt

	logrus.Println("Server is shutting down...")
	srv.Shutdown(context.Background())

	logrus.Println("Server is shut down")

	os.Exit(0)
}
