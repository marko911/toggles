package middleware

import (
	"context"
	"net/http"
	"toggle/server/pkg/auth"

	"github.com/urfave/cli/v2"
)

//Authorizer binds the initiated authorizer to context
func Authorizer(ctx *cli.Context, a *auth.Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), auth.ServiceKey, a))
			next.ServeHTTP(w, r)
		})

	}
}
