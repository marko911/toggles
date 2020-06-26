package handler

import (
	"net/http"
	"toggle/server/pkg/middleware"
	"toggle/server/pkg/models"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

// Router contains all endpoints and provides a handler
type Router struct {
	Db models.Session
}

// Handler returns an http.Handler encompassing all endpoint routes
func (r *Router) Handler(ctx *cli.Context) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/flags", FlagsHandler)
	router.Use(
		middleware.Store(ctx, r.Db),
	)
	return router
}
