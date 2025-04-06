package goergohandler

import "net/http"

type ResponseWithHttpStatus struct {
	HttpStatusCode int
	Response       any
}

// SetHeader sets the header for the response
func (r *ResponseWithHttpStatus) SetHeaders(w http.ResponseWriter) {
	if r.HttpStatusCode != 0 {
		w.WriteHeader(r.HttpStatusCode)
	}
}

func TryResponseWithHttpStatus(response any) (ResponseWithHttpStatus, bool) {
	r, ok := response.(ResponseWithHttpStatus)
	return r, ok
}
