package goergohandler

import (
	"context"
	"net/http"
)

type ValueParser interface {
	// ParseRequest parses the request and returns the context with the parsed value (if any) and error.
	// New context is not attached to the request by the parser. It happens later in the middleware.
	ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error)
}

type ParserAdder interface {
	AddParser(parser ValueParser)
}

type WithValidation interface {
	Validate() error
}

type WithParser[T any] interface {
	Parse(ctx context.Context, v string) (T, error)
}
