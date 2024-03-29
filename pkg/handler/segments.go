package handler

import (
	"net/http"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/create"
	"toggle/server/pkg/models"
	"toggle/server/pkg/read"

	"github.com/sirupsen/logrus"
)

// SegmentsHandler routes segments requests
func SegmentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		HandleSegmentsGet(w, r)
		return
	case "POST":
		HandleSegmentsPost(w, r)
		return
		// case "DELETE": TODO:implement
		// 	return
	}

	respondHTTPErr(w, r, http.StatusNotFound)

}

// HandleSegmentsGet returns all segments from db
func HandleSegmentsGet(w http.ResponseWriter, r *http.Request) {
	s := read.FromContext(r.Context())
	tenant := auth.TenantFromContext(r.Context())

	c, err := s.GetSegments(*tenant)
	if err != nil {
		logrus.Error("Getting segments failed: ", err)
		respondHTTPErr(w, r, http.StatusBadRequest)
		return
	}
	respond(w, r, http.StatusOK, c)
}

// HandleSegmentsPost adds a Segment to database
func HandleSegmentsPost(w http.ResponseWriter, r *http.Request) {
	s := create.FromContext(r.Context())
	tenant := auth.TenantFromContext(r.Context())

	segment := &models.Segment{Tenant: tenant.ID}

	if err := decodeBody(r, segment); err != nil {
		RespondErr(w, r, http.StatusBadRequest, "failed to read segment from request", err)
		return
	}

	if err := s.CreateSegment(segment); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	respond(w, r, http.StatusCreated, "Segment created successfully")
}
