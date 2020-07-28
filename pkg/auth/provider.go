package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"toggle/server/pkg/create"
	"toggle/server/pkg/models"
	"toggle/server/pkg/read"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// Jwks stores the JSON web keys from Auth0
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

//JSONWebKeys is structure of webkeys
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

//ServiceKey is used for binding Authorizer to context
const ServiceKey string = "authCTXKey"

// Authorizer implements auth middleware requirements
type Authorizer struct {
}

//Middleware is the auth middleware using auth0
type Middleware interface {
	GetHandler() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type userInfo struct {
	Email string `json:"email"`
}

// TennantMiddleware finds the tenant based on token and binds to request context
func TennantMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := jwtmiddleware.FromAuthHeader(r)
	if err != nil {
		fmt.Println("NO TOKEN zzzzzzzzzzzzzzzzz")
	}
	c := GetTenantCache()
	tenantID := c.GetByAuthToken(token)

	t := &models.Tenant{}
	if tenantID == nil {
		a := FromContext(r.Context())
		user := a.GetUserInfo(token)

		s := read.FromContext(r.Context())
		t = s.GetTenant(user)
		if t == nil {
			service := create.FromContext(r.Context())
			t, err = service.CreateTenant(user)
			if err != nil {
				logrus.Error(err.Error())
			}
			fmt.Println("addding to ctx -------------", t)
		}
		c.tenants[token] = t.ID

	} else {
		t.ID = *tenantID

	}
	r = r.WithContext(context.WithValue(r.Context(), models.TenantKey, *t))
	next(w, r)
}

// GetHandler adds Auth0 middleware that handles verifying token and pem certs and returns a handler with next
func (a *Authorizer) GetHandler() func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{

		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := os.Getenv("AUTH0_AUDIENCE")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			iss := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("unvalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	return jwtMiddleware.HandlerWithNext
}

func (a *Authorizer) GetUserInfo(token string) string {
	endpoint := "https://" + os.Getenv("AUTH0_DOMAIN") + "/userinfo"
	bearer := "Bearer " + token
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logrus.Error("Error creating request object for user info auth0")
		return ""
	}
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	u := userInfo{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		logrus.Error("Couldnt unmarshal user")
		return ""
	}

	return u.Email
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
