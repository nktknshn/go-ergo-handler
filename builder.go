package goergohandler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

type HandlerBuilder interface {
	AddParser(parser ValueParser)
}

type Builder struct {
	parsers           []ValueParser
	middlewares       []MiddlewareFunc
	handlerErrorFunc  HandleErrorFunc
	handlerResultFunc HandleResultFunc
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Success bool `json:"success"`
	Result  any  `json:"result"`
}

// by default, the error will be marshalled to json {"error": "error message"}
var DefaultHandlerErrorFunc HandleErrorFunc = func(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err := err.(type) {
	case ErrorWithHttpStatus:
		err.SetHeaders(w)
	default:
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

// by default, the result will be marshalled to json {"success": true, "result": result}
var DefaultHandlerResultFunc HandleResultFunc = func(_ context.Context, w http.ResponseWriter, _ *http.Request, result any) {
	resultData := result
	w.Header().Set("Content-Type", "application/json")

	switch result := result.(type) {
	case ResponseWithHttpStatus:
		result.SetHeaders(w)
		resultData = result.Response
	default:
		w.WriteHeader(http.StatusOK)
	}

	if resultData == nil {
		resultData = struct{}{}
	}

	bs, err := json.Marshal(successResponse{Success: true, Result: resultData})
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

func New() *Builder {
	return &Builder{}
}

func (b *Builder) AddParser(parser ValueParser) {
	b.parsers = append(b.parsers, parser)
	b.middlewares = append(b.middlewares, ValueParserToMiddleware(parser, b.handlerErrorFunc))
}

// WithHandlerErrorFunc sets a function that will be called when an error is returned by some of the middlewares
func (b *Builder) WithHandlerErrorFunc(f HandleErrorFunc) *Builder {
	b.handlerErrorFunc = f
	return b
}

// WithHandlerResultFunc sets a function that will be called when a result is returned by some of the middlewares
func (b *Builder) WithHandlerResultFunc(f HandleResultFunc) *Builder {
	b.handlerResultFunc = f
	return b
}

// BuildHandler builds a handler that will call the given function after the middlewares.
func (b *Builder) BuildHandler(f func(h http.ResponseWriter, r *http.Request)) http.Handler {
	return b.ApplyMiddleware(http.HandlerFunc(f))
}

// BuildHandlerWrapped builds a handler that is wrapped with result and error handlers.
// By default, the result will be marshalled to json {"success": true, "result": result} and the error will be marshalled to json {"error": "error message"}.
// Default failure HTTP status codes are 400 for request parsing and 500 when handler returns an error.
// Success HTTP status code is 200.
// This can be changed by setting the HandlerErrorFunc and HandlerResultFunc or by returning a ErrorWithHttpStatus/ResponseWithHttpStatus from the handler or parsers.
func (b *Builder) BuildHandlerWrapped(f func(h http.ResponseWriter, r *http.Request) (any, error)) http.Handler {

	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := f(w, r)

		if err != nil && b.handlerErrorFunc != nil {
			b.handlerErrorFunc(r.Context(), w, r, err)
			return
		}

		if err != nil {
			DefaultHandlerErrorFunc(r.Context(), w, r, err)
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

// HandleErrorFunc is a function that will be called when the handler returns an error
type HandleErrorFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error)

// HandleResultFunc is a function that will be called when the handler returns a result
type HandleResultFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, result any)

func (b *Builder) ApplyMiddleware(hh http.Handler) http.Handler {
	for i := len(b.middlewares) - 1; i >= 0; i-- {
		hh = b.middlewares[i](hh)
	}
	return hh
}

func ValueParserToMiddleware(parser ValueParser, handlerErrorFunc HandleErrorFunc) MiddlewareFunc {
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, err := parser.ParseRequest(r.Context(), w, r)
			if err != nil {
				if handlerErrorFunc != nil {
					handlerErrorFunc(ctx, w, r, err)
					return
				}
				DefaultHandlerErrorFunc(ctx, w, r, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}
