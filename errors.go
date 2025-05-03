package goergohandler

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	defaultHttpStatusCodeErrInternal = http.StatusInternalServerError
)

// builderCastError is thrown if the value in context cannot be casted to the expected value
var builderCastError = errors.New("failed to cast")

// builderMissingKeyError is thrown if context is missing a key
var builderMissingKey = errors.New("missing key")

var ErrStopPropagation = errors.New("stop")

func newBuilderCastError(msg string) error {
	return fmt.Errorf("%w: %s", builderCastError, msg)
}

func newBuilderMissingKeyError(msg string) error {
	return fmt.Errorf("%w: %s", builderMissingKey, msg)
}

// ErrorWithResponseWriter is an error that can write the response themself.
type ErrorWithResponseWriter interface {
	WriteResponse(w http.ResponseWriter)
}

type ErrorWithHeaderWriter interface {
	WriteHeader(w http.ResponseWriter)
}

// ErrorWithHttpStatus is an error that has an HTTP status code.
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

func (e ErrorWithHttpStatus) WriteHeader(w http.ResponseWriter) {
	if e.HttpStatusCode != 0 {
		w.WriteHeader(e.HttpStatusCode)
	}
}

// NewError creates a new ErrorWithHttpStatus with the given status code and wrapped error.
func NewError(status int, err error) ErrorWithHttpStatus {
	return ErrorWithHttpStatus{HttpStatusCode: status, Err: err}
}

// NewErrorStr creates a new ErrorWithHttpStatus with the given status code and wrapped error string.
func NewErrorStr(status int, errS string) ErrorWithHttpStatus {
	return ErrorWithHttpStatus{HttpStatusCode: status, Err: errors.New(errS)}
}

// WrapWithStatusCode wraps an error with an HTTP status code. If the error is already implementing the interface, it returns the original error.
func WrapWithStatusCode(err error, code int) error {
	if err == nil {
		return err
	}
	if IsWrappedError(err) {
		return err
	}
	return NewError(code, err)
}

// InternalServerError hides the original error from the client.
type InternalServerError struct {
	msg string
	Err error
}

func (e InternalServerError) Error() string {
	return e.Err.Error()
}

func (e InternalServerError) WriteResponse(w http.ResponseWriter) {
	w.WriteHeader(defaultHttpStatusCodeErrInternal)
	_, _ = w.Write([]byte(`{"error":"` + e.msg + `"}`))
}

func (e InternalServerError) Unwrap() error {
	return e.Err
}

// IsWrappedError checks if the error is already wrapped with ErrorWithHeaderWriter or ErrorWithResponseWriter.
func IsWrappedError(err error) bool {
	switch err.(type) {
	case ErrorWithHeaderWriter:
		return true
	case ErrorWithResponseWriter:
		return true
	}
	return false
}

// NewInternalServerError creates an error that does not expose the original error to the client
// while wrapping the original error.
func NewInternalServerError(err error) error {
	if err == nil {
		return nil
	}
	if IsWrappedError(err) {
		return err
	}
	return InternalServerError{Err: err, msg: "internal server error"}
}

// NewInternalServerErrorExpose creates an error that exposes the original error to the client.
func NewInternalServerErrorExpose(err error) error {
	return WrapWithStatusCode(err, defaultHttpStatusCodeErrInternal)
}
