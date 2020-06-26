package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"toggle/server/pkg/handler"
	"toggle/server/pkg/models"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/mgo.v2"
)

// Server for toggle API
type Server struct {
	port   string
	router handler.Router
	db     models.Session
}

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

// NewServer builds an http.Server with a router and port
func NewServer(c *cli.Context) *Server {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	port := c.String("server-address")
	db, err := models.NewSession(&mgo.DialInfo{
		Addrs:    []string{c.String("database-address")},
		Username: c.String("mongo-username"),
		Password: c.String("mongo-password"),
	})

	r := handler.Router{Db: db}
	if err != nil {
		logrus.Fatal(err)
		return nil
	}

	return &Server{port, r, db}
}

//Start the server
func (s Server) Start(c *cli.Context) {
	logrus.SetLevel(logrus.InfoLevel)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	srv := &http.Server{
		Addr:    s.port,
		Handler: s.router.Handler(c),
	}
	go func() {
		logrus.Fatal(srv.ListenAndServe())
	}()

	logrus.Printf("Server Started at to http://localhost%s/", s.port)

	// Wait for the interrupt signal.
	<-interrupt

	logrus.Println("Server is shutting down...")
	s.db.Close()
	srv.Shutdown(context.Background())

	logrus.Println("Server is shut down")

	os.Exit(0)
}
