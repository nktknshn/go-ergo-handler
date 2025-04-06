package goergohandler

import (
	"context"
	"net/http"
)

type QueryParamWithParserType[T WithParser[T]] struct {
	Name       string
	ErrMissing error
}

func QueryParamWithParser[T WithParser[T]](name string, errMissing error) *QueryParamWithParserType[T] {
	return &QueryParamWithParserType[T]{
		Name:       name,
		ErrMissing: errMissing,
	}
}

func (p *QueryParamWithParserType[T]) Attach(b ParserAdder) *AttachedQueryParamWithParser[T] {
	a := &AttachedQueryParamWithParser[T]{
		qp: p,
	}
	b.AddParser(a)
	return a
}

type AttachedQueryParamWithParser[T WithParser[T]] struct {
	qp *QueryParamWithParserType[T]
}

func (p *AttachedQueryParamWithParser[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if !r.URL.Query().Has(p.qp.Name) {
		return ctx, WrapError(p.qp.ErrMissing, defaultHttpStatusCodeErrParsing)
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

func (p *AttachedQueryParamWithParser[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedQueryParamWithParser[T]) GetContext(ctx context.Context) T {
	v := ctx.Value(queryParamKeyType(p.qp.Name))
	if v == nil {
		panic(newBuilderCastError("error casting..."))
	}
	return v.(T)
}
