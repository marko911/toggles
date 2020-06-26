package handler

import (
	"net/http"
	"os"
	"toggle/server/pkg/models"

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
	s := models.SessionFromContext(r.Context()).Copy()
	tenant := models.TenantFromContext(r.Context())

	defer func() {
		s.Close()
	}()

	d := s.DB(os.Getenv("DB_NAME"))
	c, err := d.GetSegments(tenant)
	if err != nil {
		logrus.Error("Getting flag failed: ", err)
		respondHTTPErr(w, r, http.StatusBadRequest)
		return
	}
	encodeBody(w, r, &c)
	respond(w, r, http.StatusOK, c)
}

// HandleSegmentsPost adds a Segment to database
func HandleSegmentsPost(w http.ResponseWriter, r *http.Request) {
	s := models.SessionFromContext(r.Context()).Copy()
	tenant := models.TenantFromContext(r.Context())

	defer s.Close()
	d := s.DB(os.Getenv("DB_NAME"))

	segment := &models.Segment{Tenant: tenant.ID}

	if err := decodeBody(r, segment); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read segment from request", err)
		return
	}

	if err := segment.Insert(d); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert segment", err)
		return
	}
	respond(w, r, http.StatusCreated, "Segment created successfully")
}
