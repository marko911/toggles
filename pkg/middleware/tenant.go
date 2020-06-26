package middleware

import (
	"context"
	"net/http"
	"toggle/server/pkg/auth"

	"github.com/urfave/cli/v2"
)

// Cache binds the cache reference to the context
func Cache(ctx *cli.Context, c *auth.Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), auth.CacheServiceKey, c))
			next.ServeHTTP(w, r)
		})
	}
}
