package handler

import (
	"fmt"
	"net/http"
	"toggle/server/pkg/models"
)

// Evaluation represents a client request for flag evaluation
type Evaluation struct {
	FlagKey string      `json:"flagKey"`
	User    models.User `json:"user"`
}

// EvaluationHandler computes the variation shown to user for given flag
func EvaluationHandler(w http.ResponseWriter, r *http.Request) {

	e, err := parseEvaluation(w, r)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to parse evaluation", err)
		return
	}
	fmt.Println("EVAL!", e.User, e.FlagKey)
	respond(w, r, http.StatusCreated, "Eval created successfully")

}

func parseEvaluation(w http.ResponseWriter, r *http.Request) (*Evaluation, error) {
	var e Evaluation

	if err := decodeBody(r, &e); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read user from request", err)
		return nil, err
	}
	return &e, nil
	//Parse user from request
	// u, ok := e["user"].(*models.User)

	// u.Tenant = tenant.ID

}
