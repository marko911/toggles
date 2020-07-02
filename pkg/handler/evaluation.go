package handler

import (
	"errors"
	"net/http"
	"sync"
	"toggle/server/pkg/create"
	"toggle/server/pkg/models"
)

// Evaluation represents a client request for flag evaluation

// EvaluationHandler computes the variation shown to user for given flag
func EvaluationHandler(w http.ResponseWriter, r *http.Request) {

	eval, err := handleEvalRequest(w, r)
	if err != nil {
		return
	}

	createService := create.FromContext(r.Context())

	errc := make(chan error, 1)
	var wg sync.WaitGroup

	// Add user to db if not existing
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = createService.CreateUser(&eval.User); err != nil {
			respondErr(w, r, http.StatusBadRequest, "failed to persist user to storage", err)
			errc <- err
		}
	}()

	// Add any custom user attributes to db
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = createService.CreateAttributes(&eval.User); err != nil {
			respondErr(w, r, http.StatusBadRequest, "failed to add custom attribute to storage", err)
			errc <- err
		}
	}()

	wg.Wait()

	select {
	case err := <-errc:
		if err != nil {
			return
		}
	default:
		break
	}

	close(errc)

	//
	respond(w, r, http.StatusCreated, "Eval created successfully")

}

func handleEvalRequest(w http.ResponseWriter, r *http.Request) (*models.Evaluation, error) {
	var e models.Evaluation

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
