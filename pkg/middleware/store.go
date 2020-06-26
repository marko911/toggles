package middleware

import (
	"context"
	"net/http"
	"toggle/server/pkg/models"

	"github.com/urfave/cli/v2"
)

// Store binds the database session to the context
func Store(ctx *cli.Context, s models.Session) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), models.MongoKey, s)) // session
			next.ServeHTTP(w, r)
		})
	}
}
