package handler

import (
	"errors"
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

	e, err := evaluationFromRequest(w, r)
	if err != nil {
		return
	}
	// TODO: start here

	fmt.Println("flagkjey! ", e.FlagKey)
	respond(w, r, http.StatusCreated, "Eval created successfully")

}

func evaluationFromRequest(w http.ResponseWriter, r *http.Request) (*Evaluation, error) {
	var e Evaluation

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
