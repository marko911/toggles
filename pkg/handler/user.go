package handler

import (
	"errors"
	"net/http"
	"os"
	"toggle/server/pkg/models"
)

// GetUser returns User by fetching or registering user from request
func GetUser(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	s := models.SessionFromContext(r.Context()).Copy()
	tenant := models.TenantFromContext(r.Context())

	d := s.DB(os.Getenv("DB_NAME"))

	var body map[string]interface{}

	if err := decodeBody(r, &body); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read user from request", err)
		return nil, err
	}

	//Parse user from request
	u, ok := body["user"].(*models.User)
	u.Tenant = tenant.ID

	if !ok {
		respondErr(w, r, http.StatusBadRequest, "failed to read user from request", &body)
		return nil, errors.New("Error")
	}

	u, err := u.UpsertUser(d)

	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert user", err)
		return nil, err
	}
	return u, body, err

}
