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

type paramBookIDWithParserType string

func (p paramBookIDWithParserType) Parse(ctx context.Context, v string) (paramBookIDWithParserType, error) {
	return paramBookIDWithParserType(v + "_parsed"), nil
}

func TestQueryParamWithParser(t *testing.T) {
	queryParam := goergohandler.QueryParamWithParser[paramBookIDWithParserType]("book_id", errors.New("book_id is required"))

	builder := goergohandler.New()
	attachedQueryParam := queryParam.Attach(builder)

	handler := builder.BuildHandler(func(w http.ResponseWriter, r *http.Request) {
		queryParam := attachedQueryParam.Get(r)
		w.Write([]byte(queryParam))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?book_id=1", nil)
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), "1_parsed")
}
