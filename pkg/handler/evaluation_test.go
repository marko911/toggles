package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"toggle/server/pkg/errors"
	"toggle/server/pkg/evaluate"
	"toggle/server/pkg/handler/testutil"
	"toggle/server/pkg/mock"
	"toggle/server/pkg/models"

	"github.com/cheekybits/is"
	"gopkg.in/mgo.v2/bson"
)

var badRequestMissingFlag = []byte(`{
	"user": {
		"key":"marko@hey.com",
		"attributes": {
			"groups": ["beta testers"],
			"gender": "male"
		}
	}
}`)

var missingFieldUser = []byte(`{
	"flagKey": "someky"
}`)

var invalidFlagKey = []byte(`{
	"flagKey":"not-existing",
	"user": {
		"key":"marko@hey.com",
		"attributes": {
			"groups": ["beta testers"],
			"gender": "male"
		}
	}
}`)

func TestEvaluationHandler(t *testing.T) {

	tests := map[string]struct {
		Body       io.Reader
		Expected   string
		CtxWrapper func(c context.Context) context.Context
	}{
		"missing flag in request":      {Body: bytes.NewBuffer(badRequestMissingFlag), Expected: errors.ErrEvalRequestMissingFlag.Error()},
		"mising user field in request": {Body: bytes.NewBuffer(missingFieldUser), Expected: errors.ErrEvalRequestMissingUser.Error()},
		"invalid flag key": {
			Body:     bytes.NewBuffer(invalidFlagKey),
			Expected: errors.ErrFlagNotFound.Error(),
			CtxWrapper: func(c context.Context) context.Context {
				mockEvalService := evaluate.NewService(&mock.EvaluateInvalidFlagKey{})
				return context.WithValue(c, evaluate.ServiceKey, mockEvalService)
			}},
	}

	ctx := mock.CreateContext(t,
		func(c context.Context) context.Context {

			tempTenant := models.Tenant{ID: bson.ObjectIdHex("5ef5f06a4fc7eb0006772c49")}
			return context.WithValue(c, models.TenantKey, tempTenant)
		},
	)
	t.Log("Given the need to be evaluate a flag:")
	{

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				is := is.New(t)

				req, err := http.NewRequest("POST", "/evaluate", tc.Body)
				is.NoErr(err)

				if tc.CtxWrapper != nil {
					ctx = tc.CtxWrapper(ctx)
				}

				req = req.WithContext(ctx)
				h, w := testutil.GetHandler(ctx, EvaluationHandler)
				h.ServeHTTP(w, req)
				respStr := w.Body.String()
				is.OK(strings.Contains(respStr, tc.Expected))
				t.Logf("\t\tShould receive a \"%s\" message. %v",
					tc.Expected, checkMark)
			})
		}

	}
}
