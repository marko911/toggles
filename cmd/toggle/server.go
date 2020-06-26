package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/create"
	"toggle/server/pkg/evaluate"
	"toggle/server/pkg/handler"
	"toggle/server/pkg/message"
	"toggle/server/pkg/read"
	"toggle/server/pkg/store/mongo"

	"github.com/bugsnag/bugsnag-go"
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
		bugsnag.Notify(err)
		logrus.Fatal(err)
		logrus.Println("ERROR STARTING MONGO")
		return nil
	}
	// db services
	mongo.PrepareDB(s)
	create := create.NewService(s)
	read := read.NewService(s)
	evaluate := evaluate.NewService(s)

	// messaging service init
	logrus.Println("Initiating nats client")
	natsClient := message.NewNatsClient(c)
	logrus.Println("Nats client initiated.")
	messenger := message.NewNatsService(natsClient)

	message.StartEvalEventsReceiever(create, messenger)
	fmt.Println("-------starting events receiver---------")
	// cache polling init
	cache := auth.GetCache()
	cache.StartPollingEvals(read)
	r := handler.Router{Create: create, Read: read, Message: messenger, Evaluate: evaluate, Authorizer: &auth.Authorizer{}, Cache: cache}

	return &Server{r, s}
}

//Start the server
func (s Server) Start(c *cli.Context) {
	defer s.store.Close()

	logrus.SetLevel(logrus.InfoLevel)
	addr := c.String("server-address")

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	srv := &http.Server{
		Addr:    addr,
		Handler: bugsnag.Handler(s.router.Handler(c)),
	}

	go func() {
		logrus.Printf("Server Started at to http://localhost%s/", addr)
		errChan <- srv.ListenAndServe()
	}()

	logrus.Fatal(<-errChan)

}
