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

type PayloadWithValidation[T PayloadWithValidationErrorType] struct {
	Payload   T
	ParserErr error
}

type PayloadWithValidationErrorType interface {
	Validate() error
}

func (p *PayloadWithValidation[T]) Attach(builder HandlerBuilder) *AttachedPayloadWithValidation[T] {
	a := &AttachedPayloadWithValidation[T]{p}
	builder.AddParser(a)
	return a
}

type AttachedPayloadWithValidation[T PayloadWithValidationErrorType] struct {
	p *PayloadWithValidation[T]
}

func (p *AttachedPayloadWithValidation[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	var pl T
	err := json.NewDecoder(r.Body).Decode(&pl)
	if err != nil {
		parseErr := ErrPayloadParserError
		if p.p.ParserErr != nil {
			parseErr = p.p.ParserErr
		}
		return ctx, WrapError(parseErr, defaultHttpStatusCodeErrParsing)
	}
	valErr := pl.Validate()
	if valErr != nil {
		return ctx, WrapError(valErr, defaultHttpStatusCodeErrParsing)
	}
	return context.WithValue(ctx, payloadKey, pl), nil
}

func (p *AttachedPayloadWithValidation[T]) GetRequest(r *http.Request) T {
	return p.Get(r.Context())
}

func (p *AttachedPayloadWithValidation[T]) Get(ctx context.Context) T {
	v := ctx.Value(payloadKey)
	if v == nil {
		panic(builderMissingKey)
	}
	return v.(T)
}

func NewPayloadWithValidation[T PayloadWithValidationErrorType](
	payload T,
) *PayloadWithValidation[T] {
	return &PayloadWithValidation[T]{
		Payload: payload,
	}
}
