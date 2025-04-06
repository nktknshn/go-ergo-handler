package goergohandler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type RouterParamWithParserType[T WithParser[T]] struct {
	Name       string
	ErrMissing error
}

func RouterParamWithParser[T WithParser[T]](name string, errMissing error) *RouterParamWithParserType[T] {
	return &RouterParamWithParserType[T]{
		Name:       name,
		ErrMissing: errMissing,
	}
}

func (rp *RouterParamWithParserType[T]) Attach(builder HandlerBuilder) *AttachedRouterParamWithParser[T] {
	a := &AttachedRouterParamWithParser[T]{rp}
	builder.AddParser(a)
	return a
}

type AttachedRouterParamWithParser[T WithParser[T]] struct {
	rp *RouterParamWithParserType[T]
}

func (p *AttachedRouterParamWithParser[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	vars := mux.Vars(r)
	v, ok := vars[p.rp.Name]
	if !ok {
		return ctx, WrapError(p.rp.ErrMissing, defaultHttpStatusCodeErrParsing)
	}
	var instance T
	vt, err := instance.Parse(ctx, v)
	if err != nil {
		return ctx, WrapError(err, defaultHttpStatusCodeErrParsing)
	}
	validatable, ok := any(vt).(WithValidation)
	if ok {
		err = validatable.Validate()
		if err != nil {
			return ctx, WrapError(err, defaultHttpStatusCodeErrParsing)
		}
	}
	return context.WithValue(ctx, routerParamKeyType(p.rp.Name), vt), nil
}

func (p *AttachedRouterParamWithParser[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedRouterParamWithParser[T]) GetContext(ctx context.Context) T {
	v := ctx.Value(routerParamKeyType(p.rp.Name))
	if v == nil {
		panic(builderMissingKey)
	}
	casted, ok := v.(T)
	if !ok {
		panic(builderCastError)
	}
	return casted
}
