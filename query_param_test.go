package goergohandler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	goergohandler "github.com/nktknshn/go-ergo-handler"
	"github.com/stretchr/testify/require"
)

func TestQueryParam_ParseRequest(t *testing.T) {
	queryParam := goergohandler.QueryParam("some_key", func(ctx context.Context, v string) (string, error) {
		return v, nil
	})

	builder := goergohandler.New()
	attachedQueryParam := queryParam.Attach(builder)

	handler := builder.BuildHandler(func(w http.ResponseWriter, r *http.Request) {
		queryParam := attachedQueryParam.Get(r)
		w.Write([]byte(queryParam))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?some_key=some_value", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), "some_value")

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), `{"error":"required query param is missing: some_key"}`)
	require.Equal(t, w.Code, http.StatusBadRequest)
}

type paramBookIDType int

func (p paramBookIDType) Validate() error {
	if p <= 0 {
		return errors.New("invalid book id")
	}
	return nil
}

func TestQueryParam_ParseRequest_WithValidation(t *testing.T) {
	queryParam := goergohandler.QueryParam("some_key", func(ctx context.Context, v string) (paramBookIDType, error) {
		vint, _ := strconv.Atoi(v)
		return paramBookIDType(vint), nil
	})

	builder := goergohandler.New()
	attachedQueryParam := queryParam.Attach(builder)

	handler := builder.BuildHandler(func(w http.ResponseWriter, r *http.Request) {
		queryParam := attachedQueryParam.Get(r)
		w.Write([]byte(strconv.Itoa(int(queryParam))))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?some_key=0", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), `{"error":"invalid book id"}`)
	require.Equal(t, w.Code, http.StatusBadRequest)
}
