package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
)

//GetHandler returns a http handler to serve
func GetHandler(ctx context.Context, h func(http.ResponseWriter, *http.Request)) (http.HandlerFunc, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(h)
	return handler, w
}
