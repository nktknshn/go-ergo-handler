package goergohandler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type MiddlewareFunc = func(http.Handler) http.Handler

// HandleErrorFunc is a function that will be called when the handler returns an error
type HandleErrorFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error)

// HandleResultFunc is a function that will be called when the handler returns a result
type HandleResultFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, result any)

// Builder is a builder for the handler.
type Builder struct {
	parsers           []ValueParser
	middlewares       []MiddlewareFunc
	handlerErrorFunc  HandleErrorFunc
	handlerResultFunc HandleResultFunc
}

func New() *Builder {
	return &Builder{}
}

// AddParser adds a parser to the builder.
// The handlerErrorFunc linked to the builder will be used to handle the error returned by the parser.
func (b *Builder) AddParser(parser ValueParser) {
	b.parsers = append(b.parsers, parser)
	b.middlewares = append(b.middlewares, ValueParserToMiddleware(parser, b.handlerErrorFunc))
}

// WithHandlerErrorFunc sets a function that will be called when an error is returned by some of the parsers
func (b *Builder) WithHandlerErrorFunc(f HandleErrorFunc) *Builder {
	b.handlerErrorFunc = f
	return b
}

// WithHandlerResultFunc sets a function that will be called when a result is returned by some of the parsers
func (b *Builder) WithHandlerResultFunc(f HandleResultFunc) *Builder {
	b.handlerResultFunc = f
	return b
}

// BuildHandler builds a handler that will call the given function after all the parsers succeed.
func (b *Builder) BuildHandler(f func(h http.ResponseWriter, r *http.Request)) http.Handler {
	return b.ApplyMiddleware(http.HandlerFunc(f))
}

// BuildHandlerWrapped builds a handler that is wrapped with result and error handlers.
// By default, the result will be marshalled to json {"result": result} and the error
// will be marshalled to json {"error": "error message"}.
// Errors that are not wrapped will be wrapped with InternalServerError
// which renders into 500 status code and plaint text message "Internal Server Error".
// Default failure HTTP status codes are 400 for request parsing and 500 for an error returned by the handler.
// Success HTTP status code is 200.
// This can be changed by setting the HandlerErrorFunc and HandlerResultFunc or by returning a ErrorWithHttpStatus/ResponseWithHttpStatus from the handler or parsers.
func (b *Builder) BuildHandlerWrapped(f func(h http.ResponseWriter, r *http.Request) (any, error)) http.Handler {
	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := f(w, r)

		if err != nil && b.handlerErrorFunc != nil {
			b.handlerErrorFunc(r.Context(), w, r, InternalServerError(err))
			return
		}

		if err != nil {
			DefaultHandlerErrorFunc(r.Context(), w, r, InternalServerError(err))
			return
		}

		if b.handlerResultFunc != nil {
			b.handlerResultFunc(r.Context(), w, r, result)
			return
		}

		DefaultHandlerResultFunc(r.Context(), w, r, result)
	})

	return b.ApplyMiddleware(wrapped)
}

// ApplyMiddleware applies the middlewares to the handler
func (b *Builder) ApplyMiddleware(hh http.Handler) http.Handler {
	for i := len(b.middlewares) - 1; i >= 0; i-- {
		hh = b.middlewares[i](hh)
	}
	return hh
}

// ValueParserToMiddleware converts a ValueParser to a MiddlewareFunc. If ParseRequest returns an error, the error will be handled by the handlerErrorFunc or DefaultHandlerErrorFunc if the handlerErrorFunc is nil.
func ValueParserToMiddleware(parser ValueParser, handlerErrorFunc HandleErrorFunc) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newctx, err := parser.ParseRequest(r.Context(), w, r)
			if err != nil {
				if handlerErrorFunc != nil {
					handlerErrorFunc(newctx, w, r, err)
					return
				}
				DefaultHandlerErrorFunc(newctx, w, r, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(newctx))
		})
	}
}

// errorResponse is the default error response to be marshalled to json {"error": "error message"}.
type errorResponse struct {
	Error string `json:"error"`
}

// successResponse is the default success response to be marshalled to json {"result": result}.
type successResponse struct {
	Result any `json:"result"`
}

// By default, the error will be marshalled to json {"error": "error message"}.
// Default http status code is 500. Return ErrorWithHttpStatus to customize the http status code.
// Implement ErrorWithResponseWriter or ErrorWithHeaderWriter for your errors to customize the response body or just headers.
// The method can be overridden by setting WithHandlerErrorFunc to builder before attaching any parsers.
var DefaultHandlerErrorFunc HandleErrorFunc = func(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {

	switch err := err.(type) {
	case ErrorWithResponseWriter:
		err.WriteResponse(w)
		return
	case ErrorWithHeaderWriter:
		w.Header().Set("Content-Type", "application/json")
		err.WriteHeader(w)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(defaultHttpStatusCodeErrInternal)
	}

	bs, err := json.Marshal(errorResponse{Error: err.Error()})
	if err != nil {
		slog.Error("error marshalling json", "error", err)
		return
	}
	_, err = w.Write(bs)
	if err != nil {
		slog.Error("error sending response", "error", err)
		return
	}
}

// By default, the result will be marshalled to json {"result": result}.
// Status code is 200. Return ResponseWithHttpStatus to customize the http status code.
// Implement ResponseWithResponseWriter for your results to customize the response body and headers.
// Nil result will be marshalled to json {"result": {}}.
// The method can be overridden by setting WithHandlerResultFunc.
var DefaultHandlerResultFunc HandleResultFunc = func(_ context.Context, w http.ResponseWriter, _ *http.Request, result any) {
	resultData := result

	switch result := result.(type) {
	case ResponseWithResponseWriter:
		result.WriteResponse(w)
		return
	case ResponseWithHttpStatus:
		w.Header().Set("Content-Type", "application/json")
		result.WriteHeaders(w)
		resultData = result.Response
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	if resultData == nil {
		resultData = struct{}{}
	}

	bs, err := json.Marshal(successResponse{Result: resultData})
	if err != nil {
		slog.Error("error marshalling json", "error", err)
		return
	}
	_, err = w.Write(bs)
	if err != nil {
		slog.Error("error sending response", "error", err)
		return
	}
}
