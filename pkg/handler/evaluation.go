package handler

import (
	"net/http"
	"sync"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/create"
	"toggle/server/pkg/errors"
	"toggle/server/pkg/evaluate"
	"toggle/server/pkg/message"
	"toggle/server/pkg/models"
	"toggle/server/pkg/read"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const messageSubject string = "evaluations"

// EvaluationHandler computes the variation shown to user for given flag
func EvaluationHandler(w http.ResponseWriter, r *http.Request) {
	read := read.FromContext(r.Context())
	// get tenant user for request
	tenant := auth.TenantFromContext(r.Context())
	if tenant == nil {
		RespondErr(w, r, http.StatusNotFound, "invalid client key")
		return
	}

	// cast request body into evaluation data
	eval := handleEvalRequest(w, r)
	if eval == nil {
		return // error occured and response error was already written to w
	}

	createService := create.FromContext(r.Context())

	var u models.User
	err := mapstructure.Decode(eval.User, &u)
	u.Tenant = tenant.ID

	errc := make(chan error, 1)
	var wg sync.WaitGroup

	// Add user to db if not existing
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := createService.CreateUser(&u); err != nil {
			RespondErr(w, r, http.StatusBadRequest, "failed to persist user to storage", err)
			errc <- err
		}
	}()

	// Add any custom user attributes to db
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := createService.CreateAttributes(&u); err != nil {
			RespondErr(w, r, http.StatusBadRequest, "failed to add custom attribute to storage", err)
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

	//222 flags, err := read.GetFlags(*tenant)

	flag, err := read.GetFlag(eval.FlagKey)
	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	// 222 var v []models.Evaluation

	// we return all flag evalutaions for this user for client sdk to have
	// avoids having to call multiple api calls at flag fork point
	// for _, flag := range flags { 222

	// TODO: invalidate cache when flag data changes
	cache := auth.CacheFromContext(r.Context())
	var matchedVariation *models.Evaluation

	if flag.HasLimit() && cache.GetEvalCount(flag.ID) > flag.Limit {
		matchedVariation, err = s.MatchDefault(*eval, flag)
		if err != nil {
			logrus.Error(err)
			RespondErr(w, r, http.StatusBadRequest, err)
			return
		}
	} else {
		matchedVariation, err = s.Evaluate(*eval, flag)
		if err != nil {
			logrus.Error(err)
			RespondErr(w, r, http.StatusBadRequest, err)
			return
		}
	}
	// v = append(v, *matchedVariation) 222

	// } 2222

	matchedVariation.User = u

	messenger := message.FromContext(r.Context())

	messenger.Publish(messageSubject, matchedVariation)

	respond(w, r, http.StatusCreated, matchedVariation)

}

func handleEvalRequest(w http.ResponseWriter, r *http.Request) *evaluate.EvaluationRequest {
	var e evaluate.EvaluationRequest

	if err := decodeBody(r, &e); err != nil {
		RespondErr(w, r, http.StatusBadRequest, "request body structure is invalid: ", err)
		return nil
	}

	if e.FlagKey == "" {
		RespondErr(w, r, http.StatusBadRequest, errors.ErrJSONPayloadInvalidFlagKey)
		return nil
	}

	var u models.User
	err := mapstructure.Decode(e.User, &u)
	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, errors.ErrEvalRequestMissingUser)
		return nil
	}

	if u.Key == "" {
		RespondErr(w, r, http.StatusBadRequest, errors.ErrEvalRequestMissingUser)
		return nil
	}

	return &e

}
