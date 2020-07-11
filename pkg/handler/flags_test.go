package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"toggle/server/pkg/mock"
	"toggle/server/pkg/models"

	"gopkg.in/mgo.v2/bson"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestHandleFlagsGet(t *testing.T) {
	ctx := mock.CreateContext(t,
		func(c context.Context) context.Context {

			tempTenant := models.Tenant{ID: bson.ObjectIdHex("5ef5f06a4fc7eb0006772c49")}
			return context.WithValue(c, models.TenantKey, tempTenant)
		},
	)

	t.Log("Given the need to be able to fetch all flags")
	{
		t.Logf("\tWhen checking \"%s\" for status code \"%d\"", "/flags", 200)
		{
			req, err := http.NewRequest("GET", "/flags", nil)
			if err != nil {
				t.Fatal(err)
			}
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(HandleFlagsGet)

			handler.ServeHTTP(w, req)

			if status := w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			t.Logf("\t\tShould receive a \"%d\" status. %v",
				w.Code, checkMark)

			content, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			if len(content) == 0 {
				t.Fatal("Home handler returned empty body")
			}
			t.Log("\t\tShould receive flags from repository.",
				checkMark)
		}

	}

}

// tests := map[string]struct {
// 	Body     io.Reader
// 	Expected string
// }{
// 	"simple":       {input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
// 	"wrong sep":    {input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
// 	"no sep":       {input: "abc", sep: "/", want: []string{"abc"}},
// 	"trailing sep": {input: "a/b/c/", sep: "/", want: []string{"a", "b", "c"}},
// }

// for name, tc := range tests {
// 	t.Run(name, func(t *testing.T) {
// 		// got := Split(tc.input, tc.sep)
// 		diff := cmp.Diff(tc.want, got)
// 		if diff != "" {
// 			t.Fatalf(diff)
// 		}
// 	})
// }
