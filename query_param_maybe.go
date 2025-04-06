package goergohandler

import (
	"context"
	"net/http"
)

type QueryParamMaybeParser[T any] func(ctx context.Context, v string) (T, error)

type QueryParamMaybe[T any] struct {
	Name   string
	Parser QueryParamMaybeParser[T]
}

func (qp *QueryParamMaybe[T]) Attach(b ParserAdder) *AttachedQueryParamMaybe[T] {
	a := &AttachedQueryParamMaybe[T]{qp}
	b.AddParser(a)
	return a
}

type AttachedQueryParamMaybe[T any] struct {
	qp *QueryParamMaybe[T]
}

func (p *AttachedQueryParamMaybe[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if !r.URL.Query().Has(p.qp.Name) {
		return ctx, nil
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
	return context.WithValue(ctx, queryParamKeyType(p.qp.Name), &v), nil
}
func (p *AttachedQueryParamMaybe[T]) GetRequestMaybe(r *http.Request) (*T, bool) {
	return p.GetMaybe(r.Context())
}

func (p *AttachedQueryParamMaybe[T]) GetMaybe(ctx context.Context) (*T, bool) {
	v := ctx.Value(queryParamKeyType(p.qp.Name))
	if v == nil {
		return nil, false
	}
	vptr, ok := v.(*T)
	if !ok {
		panic(newBuilderCastError("error casting..."))
	}
	return vptr, true
}

func NewQueryParamMaybe[T any](
	name string,
	parser QueryParamMaybeParser[T],
) *QueryParamMaybe[T] {
	return &QueryParamMaybe[T]{
		name, parser,
	}
}
