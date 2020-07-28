package handler

import (
	"net/http"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/create"
	"toggle/server/pkg/evaluate"
	"toggle/server/pkg/middleware"
	"toggle/server/pkg/models"
	"toggle/server/pkg/read"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2/bson"
)

var tempTenant models.Tenant = models.Tenant{ID: bson.ObjectIdHex("5ef5f06a4fc7eb0006772c49")}

// Router contains all endpoints and provides a handler
type Router struct {
	Create      create.Service
	Read        read.Service
	Evaluate    evaluate.Service
	Authorizer  *auth.Authorizer
	TenantCache *auth.TenantCache
}

// Handler returns an http.Handler encompassing all endpoint routes
func (r *Router) Handler(ctx *cli.Context) http.Handler {
	router := mux.NewRouter()
	tenantRoutes := mux.NewRouter().PathPrefix("/api").Subrouter()

	tenantRoutes.HandleFunc("/flags", FlagsHandler)
	tenantRoutes.HandleFunc("/segments", SegmentsHandler)

	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(r.Authorizer.GetHandler()),
		negroni.HandlerFunc(auth.TennantMiddleware),
		negroni.Wrap(tenantRoutes),
	))

	router.HandleFunc("/evaluate", EvaluationHandler).Methods("POST")

	router.Use(
		cors(ctx),
		middleware.Services(ctx, r.Create, r.Read, r.Evaluate),
		middleware.Authorizer(ctx, r.Authorizer),
		middleware.TenantCache(ctx, r.TenantCache),
	)

	return router
}
