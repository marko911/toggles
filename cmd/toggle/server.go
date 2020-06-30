package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"toggle/server/pkg/create"
	"toggle/server/pkg/handler"
	"toggle/server/pkg/read"
	"toggle/server/pkg/store/mongo"

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
	fmt.Println(e)
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
	create := create.NewService(s)
	read := read.NewService(s)

	r := handler.Router{Create: create, Read: read}
	if err != nil {
		logrus.Fatal(err)
		return nil
	}

	return &Server{r, s}
}

//Start the server
func (s Server) Start(c *cli.Context) {
	logrus.SetLevel(logrus.InfoLevel)

	addr := c.String("server-address")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

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
	s.store.Close()
	srv.Shutdown(context.Background())

	logrus.Println("Server is shut down")

	os.Exit(0)
}
