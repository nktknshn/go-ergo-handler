package goergohandler

import (
	"context"
	"net/http"
)

type ValueParser interface {
	ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error)
}

type ParserAdder interface {
	AddParser(parser ValueParser)
}
