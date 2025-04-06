package goergohandler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	goergohandler "github.com/nktknshn/go-ergo-handler"
	"github.com/stretchr/testify/require"
)

type queryParamValue string

func TestQueryParamMaybe_ParseRequest(t *testing.T) {
	queryParam := goergohandler.QueryParamMaybe("some_key", func(ctx context.Context, v string) (queryParamValue, error) {
		if v == "" {
			return "", errors.New("query param is empty")
		}
		return queryParamValue(v), nil
	})

	builder := goergohandler.New()
	attachedQueryParam := queryParam.Attach(builder)

	handler := builder.BuildHandler(func(w http.ResponseWriter, r *http.Request) {
		queryParam, _ := attachedQueryParam.GetMaybe(r)
		if queryParam == nil {
			w.Write([]byte(`NO QUERY PARAM`))
			return
		}
		w.Write([]byte(*queryParam))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?some_key=some_value", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), "some_value")

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), `NO QUERY PARAM`)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/?some_key=", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), `{"error":"query param is empty"}`)
	require.Equal(t, w.Code, http.StatusBadRequest)
}
