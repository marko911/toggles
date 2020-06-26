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

	"github.com/bugsnag/bugsnag-go"
	bugsnagnegroni "github.com/bugsnag/bugsnag-go/negroni"
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

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "Hi From default")
}

// Handler returns an http.Handler encompassing all endpoint routes
func (r *Router) Handler(ctx *cli.Context) http.Handler {
	errorHandlerConfig := bugsnag.Configuration{
		// Your Bugsnag project API key
		APIKey: ctx.String("busnag-api-key"),
		// The import paths for the Go packages containing your source files
		ProjectPackages: []string{"main", "toggles/server"},
	}

	router := mux.NewRouter()
	router.Path("/").HandlerFunc(defaultHandler)
	tenantRoutes := mux.NewRouter().PathPrefix("/api").Subrouter()
	tenantRoutes.Use(
		cors(ctx),
	)
	evalRoutes := mux.NewRouter().PathPrefix("/evals").Subrouter()

	evalRoutes = evalRoutes.StrictSlash(true)

	evalRoutes.Use(corsEvals())

	tenantRoutes.HandleFunc("/flags/{id}", FlagHandler)
	tenantRoutes.HandleFunc("/flags", FlagsHandler)
	tenantRoutes.HandleFunc("/segments", SegmentsHandler)
	tenantRoutes.HandleFunc("/segments/{id}", SegmentHandler)

	evalRoutes.HandleFunc("/flags/{clientKey}", EvaluationHandler).Methods("POST")
	evalRoutes.HandleFunc("/record/{clientKey}", RecordHandler).Methods("POST")

	evalRoutes.Use(
		auth.EvalTenantMiddleware(ctx),
	)

	n := negroni.New(
		negroni.HandlerFunc(r.Authorizer.GetHandler()),
		negroni.HandlerFunc(auth.TennantMiddleware),
		negroni.Wrap(tenantRoutes),
		bugsnagnegroni.AutoNotify(errorHandlerConfig),
	)

	router.PathPrefix("/api").Handler(n)

	router.PathPrefix("/evals").Handler(negroni.New(
		negroni.Wrap(evalRoutes),
		bugsnagnegroni.AutoNotify(errorHandlerConfig),
	))

	router.Use(
		middleware.Services(ctx, r.Create, r.Read, r.Evaluate, r.Message),
		middleware.Authorizer(ctx, r.Authorizer),
		middleware.Cache(ctx, r.Cache),
	)

	return router
}
