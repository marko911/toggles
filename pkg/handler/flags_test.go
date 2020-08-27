package handler

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"toggle/server/pkg/auth"
	"toggle/server/pkg/errors"
	"toggle/server/pkg/mock"
	"toggle/server/pkg/models"

	"github.com/cheekybits/is"
	"gopkg.in/mgo.v2/bson"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestHandleFlagsGet(t *testing.T) {
	is := is.New(t)
	ctx := mock.CreateContext(t,
		func(c context.Context) context.Context {

			tempTenant := &models.Tenant{ID: bson.ObjectIdHex("5ef5f06a4fc7eb0006772c49")}
			return context.WithValue(c, auth.TenantKey, tempTenant)
		},
	)

	t.Log("Given the need to be able to fetch all flags")
	{
		t.Logf("\tWhen checking \"%s\" for status code \"%d\" and response content", "/flags", 200)
		{
			req, err := http.NewRequest("GET", "/flags", nil)

			is.NoErr(err)

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
			is.NoErr(err)
			is.OK(len(content))
			t.Log("\t\tShould receive flags from repository.",
				checkMark)
		}

	}

}

var postBadPayloadFormat = []byte(`{
	"name": "a flag"
}`)

var postBadPayloadVariations = []byte(`{
	"name":"some flag",
	"key" :"some-flag",
	"variations":[
		{

		}
	]
}`)

var postGoodPayload = []byte(`{
	"name": "Early testers",
	"key": "beta-testers",
	"enabled": true,
	"variations": [
		{
			"name": "On",
			"percent": 100
		},
		{
			"name": "Off",
			"percent": 0
		}
	]
}`)

func TestHandleFlagsPost(t *testing.T) {
	cases := map[string]struct {
		Body     io.Reader
		Expected string
	}{
		"bad payload format":      {Body: bytes.NewBuffer(postBadPayloadFormat), Expected: errors.ErrJSONPayloadInvalidFlagKey.Error()},
		"bad payload: variations": {Body: bytes.NewBuffer(postBadPayloadVariations), Expected: errors.ErrJSONPayloadInvalidVariationEmpty.Error()},
		"good payload":            {bytes.NewBuffer(postGoodPayload), SuccessFlagCreated},
	}

	ctx := mock.CreateContext(t,
		// AUTH context injection
		func(c context.Context) context.Context {
			tempTenant := &models.Tenant{ID: bson.ObjectIdHex("5ef5f06a4fc7eb0006772c49")}
			return context.WithValue(c, auth.TenantKey, tempTenant)
		},
	)

	t.Log("Given the need to be able to create a new flag:")
	{
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				is := is.New(t)
				req, err := http.NewRequest("POST", "/flags", tc.Body)
				if err != nil {
					t.Fatal(err)
				}
				req = req.WithContext(ctx)
				w := httptest.NewRecorder()

				handler := http.HandlerFunc(HandleFlagsPost)
				handler.ServeHTTP(w, req)
				respStr := w.Body.String()
				is.OK(strings.Contains(respStr, tc.Expected))
				t.Logf("\t\tShould receive a \"%s\" message. %v",
					tc.Expected, checkMark)
			})
		}

	}
}
