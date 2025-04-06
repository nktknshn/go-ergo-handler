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

func TestQueryParam_ParseRequest(t *testing.T) {
	queryParam := goergohandler.NewQueryParam("some_key", func(ctx context.Context, v string) (string, error) {
		return v, nil
	}, errors.New("some_key is required"))

	builder := goergohandler.New()
	attachedQueryParam := queryParam.Attach(builder)

	handler := builder.BuildHandler(func(w http.ResponseWriter, r *http.Request) {
		queryParam := attachedQueryParam.GetRequest(r)
		w.Write([]byte(queryParam))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?some_key=some_value", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), "some_value")

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), `{"error":"some_key is required"}`)
	require.Equal(t, w.Code, http.StatusBadRequest)
}
