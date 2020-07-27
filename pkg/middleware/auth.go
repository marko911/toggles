package middleware

import (
	"context"
	"net/http"
	"toggle/server/pkg/auth"

	"github.com/urfave/cli/v2"
	"github.com/urfave/negroni"
)

// Auth middleware
func Auth(a auth.Middleware) negroni.Handler {
	// jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
	// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	// 		// Verify 'aud' claim
	// 		aud := os.Getenv("AUTH0_AUDIENCE")
	// 		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
	// 		if !checkAud {
	// 			return token, errors.New("Invalid audience.")
	// 		}
	// 		// Verify 'iss' claim
	// 		iss := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
	// 		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	// 		if !checkIss {
	// 			return token, errors.New("Invalid issuer.")
	// 		}

	// 		cert, err := getPemCert(token)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}

	// 		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	// 		return result, nil
	// 	},
	// 	SigningMethod: jwt.SigningMethodRS256,
	// })
	return negroni.HandlerFunc(a.GetHandler())
}

//Authorizer binds the initiated authorizer to context
func Authorizer(ctx *cli.Context, a *auth.Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), auth.ServiceKey, a))
		})

	}
}
