package goergohandler

import (
	"context"
	"net/http"
)

type QueryParamWithParserMaybeType[T WithParser[T]] struct {
	Name string
}

func QueryParamWithParserMaybe[T WithParser[T]](name string) *QueryParamWithParserMaybeType[T] {
	return &QueryParamWithParserMaybeType[T]{
		Name: name,
	}
}

func (p *QueryParamWithParserMaybeType[T]) Attach(b ParserAdder) *AttachedQueryParamWithParserMaybe[T] {
	a := &AttachedQueryParamWithParserMaybe[T]{
		qp: p,
	}
	b.AddParser(a)
	return a
}

type AttachedQueryParamWithParserMaybe[T WithParser[T]] struct {
	qp *QueryParamWithParserMaybeType[T]
}

func (p *AttachedQueryParamWithParserMaybe[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if !r.URL.Query().Has(p.qp.Name) {
		return ctx, nil
	}
	var instance T
	vstr := r.URL.Query().Get(p.qp.Name)
	v, err := instance.Parse(ctx, vstr)
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

func (p *AttachedQueryParamWithParserMaybe[T]) GetMaybe(r *http.Request) (*T, bool) {
	return p.GetContextMaybe(r.Context())
}

func (p *AttachedQueryParamWithParserMaybe[T]) GetDefault(r *http.Request, defaultVal T) T {
	v, ok := p.GetMaybe(r)
	if !ok {
		return defaultVal
	}
	return *v
}

func (p *AttachedQueryParamWithParserMaybe[T]) GetContextMaybe(ctx context.Context) (*T, bool) {
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
