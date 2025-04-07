package goergohandler

import (
	"context"
	"net/http"
)

// ParseRequest parses the request and returns the context with the parsed value (if any) and error.
// New context is not attached to the request by the parser. It happens later in the middleware.
type ValueParser interface {
	ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error)
}

// ParserAdder is an interface implemented by Builder and used by Parsers to add themselves to the Builder.
type ParserAdder interface {
	AddParser(parser ValueParser)
}

// WithValidation is an interface that can be implemented by a type to validate the parsed value.
type WithValidation interface {
	Validate() error
}

// WithParser is an interface that can be implemented by a type to parse the value from a string.
type WithParser[T any] interface {
	Parse(ctx context.Context, v string) (T, error)
}
