package handler

import (
	"errors"
	"net/http"
	"toggle/server/pkg/create"
	e "toggle/server/pkg/errors"
	"toggle/server/pkg/models"

	"github.com/gorilla/mux"
)

// SegmentHandler routes request on single segments
func SegmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		HandleSegmentPut(w, r)
		return
	}
}

// HandleSegmentPut updates segment in the db and returns updated result
func HandleSegmentPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		RespondErr(w, r, http.StatusBadRequest, errors.New("segment id required"))
		return
	}
	s := create.FromContext(r.Context())
	seg := &models.Segment{}

	if err := decodeBody(r, seg); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}

	if err := s.UpdateSegment(seg); err != nil {
		RespondErr(w, r, http.StatusBadRequest, e.ErrFailedUpdateSegment, err)
		return
	}
	respond(w, r, http.StatusCreated, seg)
}
