package goergohandler

import "net/http"

type ResponseWithResponseWriter interface {
	WriteResponse(w http.ResponseWriter)
}

type ResponseWithHeaderWriter interface {
	WriteHeaders(w http.ResponseWriter)
}

type ResponseWithHttpStatus struct {
	HttpStatusCode int
	Response       any
}

func NewResponseWithHttpStatus(httpStatusCode int, response any) ResponseWithHttpStatus {
	return ResponseWithHttpStatus{
		HttpStatusCode: httpStatusCode,
		Response:       response,
	}
}

// SetHeader sets the header for the response
func (r *ResponseWithHttpStatus) WriteHeaders(w http.ResponseWriter) {
	if r.HttpStatusCode != 0 {
		w.WriteHeader(r.HttpStatusCode)
	}
}

func TryResponseWithHttpStatus(response any) (ResponseWithHttpStatus, bool) {
	r, ok := response.(ResponseWithHttpStatus)
	return r, ok
}
