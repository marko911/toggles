package handler

import (
	"net/http"
	"toggle/server/pkg/create"
	"toggle/server/pkg/errors"
	"toggle/server/pkg/models"
	"toggle/server/pkg/read"

	"github.com/sirupsen/logrus"
)

// FlagsHandler routes flag requests
func FlagsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		HandleFlagsGet(w, r)
		return
	case "POST":
		HandleFlagsPost(w, r)
		return
		// case "DELETE": TODO: implemenet
		// 	HandleFlagsDelete(w, r)
		// 	return
	}
	respondHTTPErr(w, r, http.StatusNotFound)

}

// HandleFlagsGet returns all flags from db
func HandleFlagsGet(w http.ResponseWriter, r *http.Request) {
	s := read.FromContext(r.Context())
	tenant := models.TenantFromContext(r.Context())

	c, err := s.GetFlags(tenant)
	if err != nil {
		logrus.Error("Getting flags failed: ", err)
		respondHTTPErr(w, r, http.StatusBadRequest)
		return
	}
	respond(w, r, http.StatusOK, c)

}

// HandleFlagsPost adds a new flag to database
func HandleFlagsPost(w http.ResponseWriter, r *http.Request) {
	s := create.FromContext(r.Context())
	tenant := models.TenantFromContext(r.Context())

	flag := &models.Flag{Tenant: tenant.ID}
	if err := decodeBody(r, flag); err != nil {
		respondErr(w, r, http.StatusBadRequest, err)
		return
	}

	if _, err := flag.Validate(); err != nil {
		respondErr(w, r, http.StatusBadRequest, err)

	}

	if err := s.CreateFlag(flag); err != nil {
		respondErr(w, r, http.StatusBadRequest, errors.ErrFailedCreateFlag, err)
		return
	}
	respond(w, r, http.StatusCreated, errors.SuccessFlagCreated)

}
