package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"toggle/server/pkg/read"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

//EvalTenantMiddleware checks api key and binds the key's tenant to context
func EvalTenantMiddleware(ctx *cli.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			clientKey := vars["clientKey"]
			read := read.FromContext(r.Context())

			// get tenant user for request
			tenant := read.GetTenantFromAPIKey(clientKey)

			if tenant == nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"message": fmt.Sprint("invalid client key"),
					},
				})
				return
			}
			//TODO:add to cache lookup apiKey:tenant
			r = r.WithContext(context.WithValue(r.Context(), TenantKey, tenant))
			next.ServeHTTP(w, r)
		})
	}

}
