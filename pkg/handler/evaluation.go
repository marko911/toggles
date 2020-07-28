package handler

import (
	"net/http"
	"sync"
	"toggle/server/pkg/create"
	"toggle/server/pkg/errors"
	"toggle/server/pkg/evaluate"
	"toggle/server/pkg/models"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// Evaluation represents a client request for flag evaluation

// EvaluationHandler computes the variation shown to user for given flag
func EvaluationHandler(w http.ResponseWriter, r *http.Request) {
	eval := handleEvalRequest(w, r)
	if eval == nil {
		return // error occured and response error was already written to w
	}
	createService := create.FromContext(r.Context())
	tenant := models.TenantFromContext(r.Context())
	u, _ := eval.User.(models.User)
	u.Tenant = tenant.ID

	errc := make(chan error, 1)
	var wg sync.WaitGroup

	// Add user to db if not existing
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := createService.CreateUser(&u); err != nil {
			respondErr(w, r, http.StatusBadRequest, "failed to persist user to storage", err)
			errc <- err
		}
	}()

	// Add any custom user attributes to db
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := createService.CreateAttributes(&u); err != nil {
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

	s := evaluate.FromContext(r.Context())
	v, err := s.Evaluate(*eval)

	if err != nil {
		logrus.Error(err)
		respondErr(w, r, http.StatusBadRequest, err)
		return
	}

	respond(w, r, http.StatusCreated, v)

}

func handleEvalRequest(w http.ResponseWriter, r *http.Request) *evaluate.EvaluationData {
	var e evaluate.EvaluationData

	if err := decodeBody(r, &e); err != nil {
		respondErr(w, r, http.StatusBadRequest, "request body structure is invalid: ", err)
		return nil
	}

	if e.FlagKey == "" {
		respondErr(w, r, http.StatusBadRequest, errors.ErrEvalRequestMissingFlag)
		return nil
	}

	var u models.User
	err := mapstructure.Decode(e.User, &u)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, errors.ErrEvalRequestMissingUser)
		return nil
	}

	if u.Key == "" {
		respondErr(w, r, http.StatusBadRequest, errors.ErrEvalRequestMissingUser)
		return nil
	}

	return &e

}
