package goergohandler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	geh "github.com/nktknshn/go-ergo-handler"
	"github.com/stretchr/testify/require"
)

func TestRouterParam(t *testing.T) {
	builder := geh.New()
	routerParam := geh.RouterParam("id", func(ctx context.Context, v string) (int, error) {
		return strconv.Atoi(v)
	})
	attached := routerParam.Attach(builder)
	router := builder.BuildHandler(func(w http.ResponseWriter, r *http.Request) {
		id := attached.GetRequest(r)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(id)))
	})

	mux := mux.NewRouter()
	mux.Handle("/books/{id:[0-9]+}", router)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/books/1", nil)
	mux.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), "1")
}
