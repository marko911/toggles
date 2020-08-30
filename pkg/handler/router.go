package handler

import (
	"net/http"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/create"
	"toggle/server/pkg/evaluate"
	"toggle/server/pkg/message"
	"toggle/server/pkg/middleware"
	"toggle/server/pkg/models"
	"toggle/server/pkg/read"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"github.com/urfave/negroni"
)

var tempTenant models.Tenant = models.Tenant{ID: bson.ObjectIdHex("5ef5f06a4fc7eb0006772c49")}

// Router contains all endpoints and provides a handler
type Router struct {
	Create     create.Service
	Read       read.Service
	Evaluate   evaluate.Service
	Message    message.Service
	Authorizer *auth.Authorizer
	Cache      *auth.Cache
}

// Handler returns an http.Handler encompassing all endpoint routes
func (r *Router) Handler(ctx *cli.Context) http.Handler {
	router := mux.NewRouter()
	tenantRoutes := mux.NewRouter().PathPrefix("/api").Subrouter()
	evalRoutes := mux.NewRouter().PathPrefix("/evals").Subrouter()
	evalRoutes = evalRoutes.StrictSlash(true)
	tenantRoutes.HandleFunc("/flags/{id}", FlagHandler)
	tenantRoutes.HandleFunc("/flags", FlagsHandler)
	tenantRoutes.HandleFunc("/segments", SegmentsHandler)
	tenantRoutes.HandleFunc("/segments/{id}", SegmentHandler)

	evalRoutes.HandleFunc("/flags/{clientKey}", EvaluationHandler).Methods("POST")
	evalRoutes.HandleFunc("/record/{clientKey}", RecordHandler).Methods("POST")

	evalRoutes.Use(
		auth.EvalTenantMiddleware(ctx),
	)

	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(r.Authorizer.GetHandler()),
		negroni.HandlerFunc(auth.TennantMiddleware),
		negroni.Wrap(tenantRoutes),
	))

	router.PathPrefix("/evals").Handler(negroni.New(
		negroni.Wrap(evalRoutes),
	))

	router.Use(
		cors(ctx),
		middleware.Services(ctx, r.Create, r.Read, r.Evaluate, r.Message),
		middleware.Authorizer(ctx, r.Authorizer),
		middleware.Cache(ctx, r.Cache),
	)

	return router
}
