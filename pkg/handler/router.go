package handler

import (
	"net/http"
	"toggle/server/pkg/create"
	"toggle/server/pkg/middleware"
	"toggle/server/pkg/models"
	"toggle/server/pkg/read"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"gopkg.in/mgo.v2/bson"
)

var tempTenant models.Tenant = models.Tenant{ID: bson.ObjectIdHex("5ef5f06a4fc7eb0006772c49")}

// Router contains all endpoints and provides a handler
type Router struct {
	Create create.Service
	Read   read.Service
}

// Handler returns an http.Handler encompassing all endpoint routes
func (r *Router) Handler(ctx *cli.Context) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/flags", FlagsHandler)
	router.HandleFunc("/segments", SegmentsHandler)
	// router.HandleFunc("/evaluate", EvaluationHandler).Methods("POST")

	router.Use(
		middleware.Store(ctx, r.Create, r.Read),
		middleware.Tenant(ctx, tempTenant),
	)

	return router
}
