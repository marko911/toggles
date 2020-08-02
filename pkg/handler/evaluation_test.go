package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"toggle/server/pkg/errors"
	"toggle/server/pkg/handler/testutil"
	"toggle/server/pkg/mock"

	"github.com/cheekybits/is"
)

var missingFieldUser = []byte(`{
	"flagKey": "someky"
}`)

func TestEvaluationHandler(t *testing.T) {

	tests := map[string]struct {
		Body       io.Reader
		Expected   string
		CtxWrapper func(c context.Context) context.Context
	}{
		"mising user field in request": {Body: bytes.NewBuffer(missingFieldUser), Expected: errors.ErrEvalRequestMissingUser.Error()},
	}

	ctx := mock.CreateContext(t)
	t.Log("Given the need to be evaluate a flag:")
	{

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				is := is.New(t)

				req, err := http.NewRequest("POST", "/evaluate/abcdef", tc.Body)
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
