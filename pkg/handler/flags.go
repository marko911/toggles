package handler

import (
	"errors"
	"fmt"
	"net/http"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/create"
	e "toggle/server/pkg/errors"
	"toggle/server/pkg/models"
	"toggle/server/pkg/read"

	"github.com/gemcook/pagination-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// FlagsResponse is the json payload for frontend
type FlagsResponse struct {
	Flags  []models.Flag `json:"flags"`
	APIKey string        `json:"apiKey"`
}

// FlagsHandler routes flag requests
func FlagsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		HandleFlagsGet(w, r)
		return
	case "POST":
		HandleFlagsPost(w, r)
		return

	}
	respondHTTPErr(w, r, http.StatusNotFound)

}

// FlagHandler is for individual flag updates
func FlagHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		HandleFlagPut(w, r)
		return
	case "GET":
		HandleFlagGet(w, r)
	}

}

// HandleFlagPut updates a flag
func HandleFlagPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		RespondErr(w, r, http.StatusBadRequest, errors.New("flag id required"))
		return
	}

	s := create.FromContext(r.Context())
	flag := &models.Flag{}

	if err := decodeBody(r, flag); err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	fmt.Println("flag", flag.Targets)
	if err := s.UpdateFlag(flag); err != nil {
		RespondErr(w, r, http.StatusBadRequest, e.ErrFailedCreateFlag, err)
		return
	}
	respond(w, r, http.StatusCreated, flag)
}

// HandleFlagGet retrieves flag evaluation records
func HandleFlagGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		RespondErr(w, r, http.StatusBadRequest, errors.New("flag id required"))
		return
	}

	p := pagination.ParseQuery(r.URL.RequestURI())

	if p.Page < 1 {
		RespondErr(w, r, http.StatusBadRequest, errors.New("Bad flag pagination request"))
		return
	}

	rs := read.FromContext(r.Context())
	flagID := bson.ObjectIdHex(id)

	evaluations, c, err := rs.GetFlagEvals(flagID, p.Page, p.Limit)

	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
	}
	stats, err := rs.GetFlagStats(flagID)
	if err != nil {
		RespondErr(w, r, http.StatusBadRequest, err)
	}

	response := map[string]interface{}{
		"id":          id,
		"total":       c,
		"page":        p.Page,
		"evaluations": evaluations,
		"stats":       stats,
	}
	respond(w, r, http.StatusAccepted, response)

}

// HandleFlagsGet returns all flags from db
func HandleFlagsGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	s := read.FromContext(r.Context())
	tenant := auth.TenantFromContext(r.Context())
	c, err := s.GetFlags(*tenant)
	if err != nil {
		logrus.Error("Getting flags failed: ", err)
		RespondErr(w, r, http.StatusBadRequest, err)
		return
	}
	respond(w, r, http.StatusOK, FlagsResponse{c, tenant.APIKEY})

}

//SuccessFlagCreated is message returned on success of flag post
var SuccessFlagCreated = "Flag created successfully"

// SuccessFlagUpdated is message for updated flag
var SuccessFlagUpdated = "Flag updated successfully"

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
		RespondErr(w, r, http.StatusBadRequest, e.ErrFailedCreateFlag, err)
		return
	}

	respond(w, r, http.StatusCreated, SuccessFlagCreated)

}
