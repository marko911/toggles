package middleware

import (
	"context"
	"net/http"
	"toggle/server/pkg/create"
	"toggle/server/pkg/evaluate"
	"toggle/server/pkg/message"
	"toggle/server/pkg/read"

	"github.com/urfave/cli/v2"
)

// Services binds all services to the context
func Services(ctx *cli.Context, c create.Service, rs read.Service, es evaluate.Service, ms message.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), create.ServiceKey, c))    // inject create service
			r = r.WithContext(context.WithValue(r.Context(), read.ServiceKey, rs))     // inject read service
			r = r.WithContext(context.WithValue(r.Context(), evaluate.ServiceKey, es)) // inject read service
			r = r.WithContext(context.WithValue(r.Context(), message.ServiceKey, ms))  // inject message service
			next.ServeHTTP(w, r)
		})
	}
}
