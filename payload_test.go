package goergohandler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	goergohandler "github.com/nktknshn/go-ergo-handler"
	"github.com/stretchr/testify/require"
)

type testPayload struct {
	SomeKey string `json:"some_key"`
}

func (p testPayload) Validate() error {
	if p.SomeKey == "" {
		return errors.New("some_key is required")
	}
	return nil
}

func TestPayloadWithValidation(t *testing.T) {

	payload := goergohandler.Payload[testPayload]()
	builder := goergohandler.New()
	attachedPayload := payload.Attach(builder)

	handler := builder.BuildHandler(func(w http.ResponseWriter, r *http.Request) {
		payload := attachedPayload.GetRequest(r)
		w.Write([]byte(payload.SomeKey))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"some_key": "some_value"}`))
	handler.ServeHTTP(w, r)

	require.Equal(t, w.Body.String(), "some_value")

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/", strings.NewReader(`{}`))
	handler.ServeHTTP(w, r)

	require.Equal(t, `{"error":"some_key is required"}`, w.Body.String())
	require.Equal(t, http.StatusBadRequest, w.Code)
}
