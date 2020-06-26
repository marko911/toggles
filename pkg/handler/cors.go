package handler

import (
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

func cors(ctx *cli.Context) func(http.Handler) http.Handler {
	allowedOrigins := ctx.StringSlice("server-allowed-hosts")

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()

			headers.Add("Access-Control-Allow-Methods", "POST, GET, PUT")
			headers.Add("Access-Control-Allow-Credentials", "true")
			// Required to be lower case for ie
			headers.Add("Access-Control-Allow-Headers", "authorization, accept, content-type, cache-control, x-feature-toggles")

			origin := strings.ToLower(r.Header.Get("Origin"))

			for _, o := range allowedOrigins {
				if origin == o {
					headers.Add("Access-Control-Allow-Origin", origin)
				}
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func corsEvals() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()

			headers.Add("Access-Control-Allow-Methods", "POST, GET")
			headers.Add("Access-Control-Allow-Credentials", "true")
			// Required to be lower case for ie
			headers.Add("Access-Control-Allow-Headers", "authorization, accept, content-type, cache-control, x-feature-toggles")

			headers.Add("Access-Control-Allow-Origin", "*")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
