package goergohandler

import (
	"context"
	"net/http"
)

const defaultHttpStatusCodeErrParsing = http.StatusBadRequest

type queryParamKeyType string

type QueryParamParserFunc[T any] func(ctx context.Context, v string) (T, error)

type QueryParam[T any] struct {
	Name       string
	Parser     QueryParamParserFunc[T]
	ErrMissing error
}

func (qp *QueryParam[T]) Attach(b ParserAdder) *AttachedQueryParam[T] {
	a := &AttachedQueryParam[T]{qp}
	b.AddParser(a)
	return a
}

type AttachedQueryParam[T any] struct {
	qp *QueryParam[T]
}

func (p *AttachedQueryParam[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if !r.URL.Query().Has(p.qp.Name) {
		return ctx, WrapError(p.qp.ErrMissing, defaultHttpStatusCodeErrParsing)
	}
	vstr := r.URL.Query().Get(p.qp.Name)
	v, err := p.qp.Parser(ctx, vstr)
	if err != nil {
		return ctx, WrapError(err, defaultHttpStatusCodeErrParsing)
	}
	validatable, ok := any(v).(WithValidation)
	if ok {
		err = validatable.Validate()
		if err != nil {
			return ctx, WrapError(err, defaultHttpStatusCodeErrParsing)
		}
	}
	return context.WithValue(ctx, queryParamKeyType(p.qp.Name), v), nil
}

func (p *AttachedQueryParam[T]) GetRequest(r *http.Request) T {
	return p.Get(r.Context())
}

func (p *AttachedQueryParam[T]) Get(ctx context.Context) T {
	v := ctx.Value(queryParamKeyType(p.qp.Name))
	if v == nil {
		// return *new(T), p.qp.ErrMissing
		panic(builderMissingKey)
	}
	return v.(T)
}

func NewQueryParam[T any](
	name string,
	parser QueryParamParserFunc[T],
	errMissing error,
) *QueryParam[T] {
	return &QueryParam[T]{
		name, parser, errMissing,
	}
}
