package middleware

import (
	"context"
	"net/http"
	"toggle/server/pkg/create"
	"toggle/server/pkg/read"

	"github.com/urfave/cli/v2"
)

// Store binds the database session to the context
func Store(ctx *cli.Context, c create.Service, r read.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), create.ServiceKey, c)) // inject create service
			r = r.WithContext(context.WithValue(r.Context(), read.ServiceKey, r))   // inject read service
			next.ServeHTTP(w, r)
		})
	}
}
