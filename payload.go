package goergohandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	defaultHttpStatusCodeErrPayloadParsing    = http.StatusBadRequest
	defaultHttpStatusCodeErrPayloadValidation = http.StatusBadRequest
)

type payloadKeyType string

var (
	payloadKey        payloadKeyType = "payload"
	ErrPayloadParsing error          = errors.New("error parsing payload")
)

type PayloadParserType[T any] struct {
	// override the default parsing error
	// TODO: implement constructor for this
	ParserErr error
}

// Payload is a parser that parses the payload from the request.
// If payload type implements WithValidation, it will be validated.
func Payload[T any]() *PayloadParserType[T] {
	return &PayloadParserType[T]{}
}

func (p *PayloadParserType[T]) Attach(builder ParserAdder) *AttachedPayloadParser[T] {
	a := &AttachedPayloadParser[T]{p}
	builder.AddParser(a)
	return a
}

// WithParsingError sets the error to be returned if the payload is not valid json.
func (p *PayloadParserType[T]) WithParsingError(err error) *PayloadParserType[T] {
	p.ParserErr = err
	return p
}

type AttachedPayloadParser[T any] struct {
	pp *PayloadParserType[T]
}

func (p *AttachedPayloadParser[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	var pl T
	err := json.NewDecoder(r.Body).Decode(&pl)
	if err != nil {
		parseErr := ErrPayloadParsing
		if p.pp.ParserErr != nil {
			parseErr = p.pp.ParserErr
		}
		return ctx, WrapWithStatusCode(parseErr, defaultHttpStatusCodeErrPayloadParsing)
	}
	err = ValidateWithValidation(pl)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrPayloadValidation)
	}
	return context.WithValue(ctx, payloadKey, pl), nil
}

func (p *AttachedPayloadParser[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedPayloadParser[T]) GetContext(ctx context.Context) T {
	return GetFromContext[T](ctx, payloadKey)
}
