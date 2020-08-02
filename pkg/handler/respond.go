package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"toggle/server/pkg/errors"
)

func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		return errors.ErrJSONPayloadEmpty
	}
	err = json.Unmarshal(body, &v)
	return err
}

func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func respond(w http.ResponseWriter, r *http.Request,
	status int, data interface{},
) {
	w.WriteHeader(status)
	if data != nil {
		encodeBody(w, r, data)
	}
}

func RespondErr(w http.ResponseWriter, r *http.Request,
	status int, args ...interface{},
) {
	respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}

func respondHTTPErr(w http.ResponseWriter, r *http.Request,
	status int,
) {
	RespondErr(w, r, status, http.StatusText(status))
}
