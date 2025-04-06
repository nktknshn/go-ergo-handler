package goergohandler_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	geh "github.com/nktknshn/go-ergo-handler"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildHandlerWrapped(t *testing.T) {
	type testCase struct {
		customErrorFunc  geh.HandleErrorFunc
		customResultFunc geh.HandleResultFunc
		name             string
		result           any
		error            error
		expectedCode     int
		expectedBody     string
		customCheck      func(t *testing.T, w *httptest.ResponseRecorder)
	}

	cases := []testCase{
		{
			name:         "success",
			result:       map[string]string{"some_key": "some_value"},
			error:        nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"success":true,"result":{"some_key":"some_value"}}`,
		},
		{
			name:         "error",
			result:       nil,
			error:        errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"some error"}`,
		},
		{
			name:         "error with code",
			result:       nil,
			error:        geh.NewError(http.StatusBadRequest, errors.New("some error")),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"some error"}`,
		},
		{
			name:         "result with code",
			result:       geh.NewResponseWithHttpStatus(http.StatusAccepted, map[string]string{"some_key": "some_value"}),
			error:        nil,
			expectedCode: http.StatusAccepted,
			expectedBody: `{"success":true,"result":{"some_key":"some_value"}}`,
		},
		{
			name: "custom error func",
			customErrorFunc: func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte(`ERROR`))
			},
			error:        errors.New("some error"),
			expectedCode: http.StatusBadRequest,
			expectedBody: `ERROR`,
			customCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, w.Header().Get("Content-Type"), "text/plain")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := geh.New()
			if c.customErrorFunc != nil {
				b.WithHandlerErrorFunc(c.customErrorFunc)
			}
			if c.customResultFunc != nil {
				b.WithHandlerResultFunc(c.customResultFunc)
			}
			handler := b.BuildHandlerWrapped(func(h http.ResponseWriter, r *http.Request) (any, error) {
				return c.result, c.error
			})

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
			require.Equal(t, c.expectedCode, w.Code)
			body, err := io.ReadAll(w.Body)
			require.NoError(t, err)
			require.Equal(t, c.expectedBody, string(body))
			if c.customCheck != nil {
				c.customCheck(t, w)
			}
		})
	}
}
