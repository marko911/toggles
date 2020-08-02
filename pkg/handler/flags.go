package handler

import (
	"fmt"
	"net/http"
	"toggle/server/pkg/auth"
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
	tenant := auth.TenantFromContext(r.Context())
	fmt.Println("tenanttttt", tenant)
	c, err := s.GetFlags(*tenant)
	if err != nil {
		logrus.Error("Getting flags failed: ", err)
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	respond(w, r, http.StatusOK, c)

}

// HandleFlagsPost adds a new flag to database
func HandleFlagsPost(w http.ResponseWriter, r *http.Request) {
	s := create.FromContext(r.Context())
	tenant := auth.TenantFromContext(r.Context())
	flag := &models.Flag{Tenant: tenant.ID}
	if err := decodeBody(r, flag); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}

	if _, err := flag.Validate(); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}

	if err := s.CreateFlag(flag); err != nil {
		RespondErr(w, r, http.StatusBadRequest, errors.ErrFailedCreateFlag, err)
		return
	}

	respond(w, r, http.StatusCreated, errors.SuccessFlagCreated)

}
