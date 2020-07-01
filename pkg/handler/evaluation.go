package handler

import (
	"errors"
	"fmt"
	"net/http"
	"toggle/server/pkg/create"
	"toggle/server/pkg/models"
)

// Evaluation represents a client request for flag evaluation
type Evaluation struct {
	FlagKey string      `json:"flagKey"`
	User    models.User `json:"user"`
}

// EvaluationHandler computes the variation shown to user for given flag
func EvaluationHandler(w http.ResponseWriter, r *http.Request) {

	eval, err := handleEvalRequest(w, r)
	if err != nil {
		return
	}

	createService := create.FromContext(r.Context())

	// Add user to db if not existing
	if err = createService.CreateUser(&eval.User); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to persist user to storage", err)
	}

	// Add any custom user attributes to db
	if err = createService.CreateAttributes(&eval.User); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to add custom attribute to storage", err)
	}

	fmt.Println("flagkjey! ", eval.FlagKey, eval.User)
	respond(w, r, http.StatusCreated, "Eval created successfully")

}

func handleEvalRequest(w http.ResponseWriter, r *http.Request) (*Evaluation, error) {
	var e Evaluation

	tenant := models.TenantFromContext(r.Context())
	e.User.Tenant = tenant.ID

	if err := decodeBody(r, &e); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read user from request", err)
		return nil, err
	}

	if e.FlagKey == "" {
		respondErr(w, r, http.StatusBadRequest, "Must provide flag key")
		return nil, errors.New("")
	}

	if e.User.Key == "" {
		respondErr(w, r, http.StatusBadRequest, "Unique user key is missing")
		return nil, errors.New("")
	}

	return &e, nil

}
