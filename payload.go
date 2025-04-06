package goergohandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type payloadKeyType string

var (
	payloadKey            payloadKeyType = "payload"
	ErrPayloadParserError error          = errors.New("payload parser error")
)

type PayloadParserType[T any] struct {
	Payload   T
	ParserErr error
}

func (p *PayloadParserType[T]) Attach(builder HandlerBuilder) *AttachedPayloadParser[T] {
	a := &AttachedPayloadParser[T]{p}
	builder.AddParser(a)
	return a
}

type AttachedPayloadParser[T any] struct {
	p *PayloadParserType[T]
}

func (p *AttachedPayloadParser[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	var pl T
	err := json.NewDecoder(r.Body).Decode(&pl)
	if err != nil {
		parseErr := ErrPayloadParserError
		if p.p.ParserErr != nil {
			parseErr = p.p.ParserErr
		}
		return ctx, WrapError(parseErr, defaultHttpStatusCodeErrParsing)
	}
	validatable, ok := any(pl).(WithValidation)
	if ok {
		err = validatable.Validate()
		if err != nil {
			return ctx, WrapError(err, defaultHttpStatusCodeErrParsing)
		}
	}
	return context.WithValue(ctx, payloadKey, pl), nil
}

func (p *AttachedPayloadParser[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedPayloadParser[T]) GetContext(ctx context.Context) T {
	v := ctx.Value(payloadKey)
	if v == nil {
		panic(builderMissingKey)
	}
	return v.(T)
}

func Payload[T any]() *PayloadParserType[T] {
	return &PayloadParserType[T]{}
}
