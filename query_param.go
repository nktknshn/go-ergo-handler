package goergohandler

import (
	"context"
	"net/http"
)

const defaultHttpStatusCodeErrParsing = http.StatusBadRequest

type queryParamKeyType string

type QueryParamParserFunc[T any] func(ctx context.Context, v string) (T, error)

type QueryParamType[T any] struct {
	Name       string
	Parser     QueryParamParserFunc[T]
	ErrMissing error
}

func (qp *QueryParamType[T]) Attach(b ParserAdder) *AttachedQueryParam[T] {
	a := &AttachedQueryParam[T]{qp}
	b.AddParser(a)
	return a
}

type AttachedQueryParam[T any] struct {
	qp *QueryParamType[T]
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

func (p *AttachedQueryParam[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedQueryParam[T]) GetContext(ctx context.Context) T {
	v := ctx.Value(queryParamKeyType(p.qp.Name))
	if v == nil {
		// return *new(T), p.qp.ErrMissing
		panic(builderMissingKey)
	}
	return v.(T)
}

func QueryParam[T any](
	name string,
	parser QueryParamParserFunc[T],
	errMissing error,
) *QueryParamType[T] {
	return &QueryParamType[T]{
		name, parser, errMissing,
	}
}
