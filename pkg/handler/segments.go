package handler

import (
	"net/http"
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
	}

	respondHTTPErr(w, r, http.StatusNotFound)

}

// HandleSegmentsGet returns all segments from db
func HandleSegmentsGet(w http.ResponseWriter, r *http.Request) {
	s := read.FromContext(r.Context())
	tenant := models.TenantFromContext(r.Context())

	c, err := s.GetSegments(tenant)
	if err != nil {
		logrus.Error("Getting segments failed: ", err)
		respondHTTPErr(w, r, http.StatusBadRequest)
		return
	}
	encodeBody(w, r, &c)
	respond(w, r, http.StatusOK, c)
}

// HandleSegmentsPost adds a Segment to database
func HandleSegmentsPost(w http.ResponseWriter, r *http.Request) {
	s := create.FromContext(r.Context())
	tenant := models.TenantFromContext(r.Context())

	segment := &models.Segment{Tenant: tenant.ID}

	if err := decodeBody(r, segment); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read segment from request", err)
		return
	}

	if err := s.CreateSegment(segment); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert segment", err)
		return
	}
	respond(w, r, http.StatusCreated, "Segment created successfully")
}
