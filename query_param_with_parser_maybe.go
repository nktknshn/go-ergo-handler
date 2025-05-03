package goergohandler

import (
	"context"
	"net/http"
)

type QueryParamWithParserMaybeType[T WithParser[T]] struct {
	Name string
}

func QueryParamWithParserMaybe[T WithParser[T]](name string) *QueryParamWithParserMaybeType[T] {
	return &QueryParamWithParserMaybeType[T]{Name: name}
}

func (p *QueryParamWithParserMaybeType[T]) Attach(b ParserAdder) *AttachedQueryParamWithParserMaybe[T] {
	a := &AttachedQueryParamWithParserMaybe[T]{p}
	b.AddParser(a)
	return a
}

type AttachedQueryParamWithParserMaybe[T WithParser[T]] struct {
	qp *QueryParamWithParserMaybeType[T]
}

func (a *AttachedQueryParamWithParserMaybe[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if !r.URL.Query().Has(a.qp.Name) {
		return ctx, nil
	}
	var instance T
	vstr := r.URL.Query().Get(a.qp.Name)
	v, err := instance.Parse(ctx, vstr)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrQueryParamParsing)
	}
	err = ValidateWithValidation(v)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrQueryParamValidation)
	}
	return context.WithValue(ctx, queryParamKeyType(a.qp.Name), v), nil
}

func (a *AttachedQueryParamWithParserMaybe[T]) GetMaybe(r *http.Request) (*T, bool) {
	return a.GetContextMaybe(r.Context())
}

func (a *AttachedQueryParamWithParserMaybe[T]) GetDefault(r *http.Request, defaultVal T) T {
	return a.GetContextDefault(r.Context(), defaultVal)
}

func (p *AttachedQueryParamWithParserMaybe[T]) GetContextMaybe(ctx context.Context) (*T, bool) {
	return GetFromContextMaybe[T](ctx, queryParamKeyType(p.qp.Name))
}

func (a *AttachedQueryParamWithParserMaybe[T]) GetContextDefault(ctx context.Context, defaultVal T) T {
	v, ok := a.GetContextMaybe(ctx)
	if !ok {
		return defaultVal
	}
	return *v
}
