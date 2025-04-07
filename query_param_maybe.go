package goergohandler

import (
	"context"
	"net/http"
)

type QueryParamMaybeParser[T any] = func(ctx context.Context, v string) (T, error)

type QueryParamMaybeType[T any] struct {
	Name   string
	Parser QueryParamMaybeParser[T]
}

// QueryParamMaybe is same as QueryParam but it doesn't return an error if the query param is missing.
// It stores a pointer to the value and returns nil if the query param is missing.
func QueryParamMaybe[T any](
	name string,
	parser QueryParamMaybeParser[T],
) *QueryParamMaybeType[T] {
	return &QueryParamMaybeType[T]{
		name, parser,
	}
}

func (qp *QueryParamMaybeType[T]) Attach(b ParserAdder) *AttachedQueryParamMaybe[T] {
	a := &AttachedQueryParamMaybe[T]{qp}
	b.AddParser(a)
	return a
}

type AttachedQueryParamMaybe[T any] struct {
	qp *QueryParamMaybeType[T]
}

func (p *AttachedQueryParamMaybe[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if !r.URL.Query().Has(p.qp.Name) {
		return ctx, nil
	}
	vstr := r.URL.Query().Get(p.qp.Name)
	v, err := p.qp.Parser(ctx, vstr)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrQueryParamParsing)
	}
	err = ValidateWithValidation(v)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrQueryParamValidation)
	}
	return context.WithValue(ctx, queryParamKeyType(p.qp.Name), &v), nil
}

func (p *AttachedQueryParamMaybe[T]) GetMaybe(r *http.Request) (*T, bool) {
	return p.GetContextMaybe(r.Context())
}

func (p *AttachedQueryParamMaybe[T]) GetDefault(r *http.Request, defaultVal T) T {
	return p.GetContextDefault(r.Context(), defaultVal)
}

func (p *AttachedQueryParamMaybe[T]) GetContextDefault(ctx context.Context, defaultVal T) T {
	v, ok := p.GetContextMaybe(ctx)
	if !ok {
		return defaultVal
	}
	return *v
}

func (p *AttachedQueryParamMaybe[T]) GetContextMaybe(ctx context.Context) (*T, bool) {
	return GetFromContextMaybe[*T](ctx, queryParamKeyType(p.qp.Name))
}
