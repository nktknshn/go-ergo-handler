package goergohandler

import (
	"context"
	"fmt"
	"net/http"
)

type QueryParamWithParserType[T WithParser[T]] struct {
	Name       string
	ErrMissing error
}

// QueryParamWithParser is same as QueryParam but it uses a parser function from the given type.
func QueryParamWithParser[T WithParser[T]](name string) *QueryParamWithParserType[T] {
	return &QueryParamWithParserType[T]{
		Name: name,
	}
}

func (p *QueryParamWithParserType[T]) WithErrMissing(errMissing error) *QueryParamWithParserType[T] {
	p.ErrMissing = errMissing
	return p
}

func (p *QueryParamWithParserType[T]) Attach(b ParserAdder) *AttachedQueryParamWithParser[T] {
	a := &AttachedQueryParamWithParser[T]{p}
	b.AddParser(a)
	return a
}

type AttachedQueryParamWithParser[T WithParser[T]] struct {
	qp *QueryParamWithParserType[T]
}

func (p *AttachedQueryParamWithParser[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if !r.URL.Query().Has(p.qp.Name) {
		err := p.qp.ErrMissing
		if err == nil {
			err = fmt.Errorf("%w: %s", ErrQueryParamMissing, p.qp.Name)
		}
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrQueryParamMissing)
	}
	var instance T
	vstr := r.URL.Query().Get(p.qp.Name)
	v, err := instance.Parse(ctx, vstr)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrQueryParamParsing)
	}
	err = ValidateWithValidation(v)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrQueryParamValidation)
	}
	return context.WithValue(ctx, queryParamKeyType(p.qp.Name), v), nil
}

func (p *AttachedQueryParamWithParser[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedQueryParamWithParser[T]) GetContext(ctx context.Context) T {
	return GetFromContext[T](ctx, queryParamKeyType(p.qp.Name))
}
