package middleware

import (
	"context"
	"net/http"
	"toggle/server/pkg/models"

	"github.com/urfave/cli/v2"
)

// Tenant binds the current tenant to the context
func Tenant(ctx *cli.Context, t models.Tenant) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), models.TenantKey, t)) // session
			next.ServeHTTP(w, r)
		})
	}
}
