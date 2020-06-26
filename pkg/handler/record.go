package handler

import (
	"net/http"
	"toggle/server/pkg/message"
	"toggle/server/pkg/models"
)

// RecordHandler receives flag evaluation call events and publishes it to the NATS topic
func RecordHandler(w http.ResponseWriter, r *http.Request) {
	messenger := message.FromContext(r.Context())
	var e models.Evaluation
	if err := decodeBody(r, &e); err != nil {
	}
	messenger.Publish("evaluations", e)
}
