package handler

import (
	"net/http"
	"os"
	"toggle/server/pkg/models"

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
		// case "DELETE":
		// 	HandleFlagsDelete(w, r)
		// 	return
	}
	respondHTTPErr(w, r, http.StatusNotFound)

}

// HandleFlagsGet returns all flags from db
func HandleFlagsGet(w http.ResponseWriter, r *http.Request) {
	s := models.SessionFromContext(r.Context()).Copy()
	tenant := models.TenantFromContext(r.Context())

	defer s.Close()

	d := s.DB(os.Getenv("DB_NAME"))

	c, err := d.GetFlags(tenant)
	if err != nil {
		logrus.Error("Getting flag failed: ", err)
		respondHTTPErr(w, r, http.StatusBadRequest)
		return
	}
	encodeBody(w, r, &c)
	respond(w, r, http.StatusOK, c)

}

// HandleFlagsPost adds a new flag to database
func HandleFlagsPost(w http.ResponseWriter, r *http.Request) {
	s := models.SessionFromContext(r.Context()).Copy()
	tenant := models.TenantFromContext(r.Context())

	defer s.Close()

	d := s.DB(os.Getenv("DB_NAME"))

	flag := &models.Flag{Tenant: tenant.ID}

	if err := decodeBody(r, flag); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read flag from request ", err)
		return
	}

	if err := flag.Insert(d); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert flag", err)
		return
	}
	respond(w, r, http.StatusCreated, "Flag created successfully")

}
