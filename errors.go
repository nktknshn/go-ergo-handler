package goergohandler

import (
	"errors"
	"fmt"
	"net/http"
)

// builderCastError происходит, когда не удалось привести значение из D к нужному типу.
// Исключительная ситуация, которая не должна происходить.
var builderCastError = errors.New("failed to cast")

// builderMissingKeyError происходит, когда не удалось получить значение из D по ключу.
// Исключительная ситуация, которая произойдет, если в Builder не добавлен нужный парсер.
var builderMissingKey = errors.New("missing key")

func newBuilderCastError(msg string) error {
	return fmt.Errorf("%w: %s", builderCastError, msg)
}

func newBuilderMissingKeyError(msg string) error {
	return fmt.Errorf("%w: %s", builderMissingKey, msg)
}

type ErrorWithHttpStatus struct {
	HttpStatusCode int
	Err            error
}

func (e ErrorWithHttpStatus) Error() string {
	return e.Err.Error()
}

func (e ErrorWithHttpStatus) Unwrap() error {
	return e.Err
}

func (e ErrorWithHttpStatus) SetHeaders(w http.ResponseWriter) {
	if e.HttpStatusCode != 0 {
		w.WriteHeader(e.HttpStatusCode)
	}
}

func NewError(status int, err error) ErrorWithHttpStatus {
	return ErrorWithHttpStatus{HttpStatusCode: status, Err: err}
}

func NewErrorStr(status int, errS string) ErrorWithHttpStatus {
	return ErrorWithHttpStatus{HttpStatusCode: status, Err: errors.New(errS)}
}

func WrapError(err error, code int) error {
	if err == nil {
		return err
	}
	_, ok := err.(ErrorWithHttpStatus)
	if ok {
		return err
	}
	return NewError(code, err)
}

func TryErrorWithHttpStatus(err error) (ErrorWithHttpStatus, bool) {
	e, ok := err.(ErrorWithHttpStatus)
	return e, ok
}
