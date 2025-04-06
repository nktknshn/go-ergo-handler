package goergohandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type routerParamKeyType string

type RouteParamParserFunc[T any] func(ctx context.Context, v string) (T, error)

type RouterParam[T any] struct {
	Name       string
	Parser     RouteParamParserFunc[T]
	ErrParsing error
	ErrMissing error
}

var ErrorRouterParamMissing = errors.New("router param missing")

func NewRouterParamMissingError(paramName string) error {
	return fmt.Errorf("%w: %s", ErrorRouterParamMissing, paramName)
}

func (rp *RouterParam[T]) Attach(builder HandlerBuilder) *AttachedRouterParam[T] {
	a := &AttachedRouterParam[T]{rp}
	builder.AddParser(a)
	return a
}

type AttachedRouterParam[T any] struct {
	rp *RouterParam[T]
}

func (p *AttachedRouterParam[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	vars := mux.Vars(r)
	v, ok := vars[p.rp.Name]
	if !ok {
		err := NewRouterParamMissingError(p.rp.Name)
		return ctx, WrapError(err, defaultHttpStatusCodeErrParsing)
	}
	vt, err := p.rp.Parser(ctx, v)
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

func (p *AttachedRouterParam[T]) GetRequest(r *http.Request) T {
	return p.Get(r.Context())
}

func (p *AttachedRouterParam[T]) Get(ctx context.Context) T {
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

func NewRouterParam[T any](name string, parser RouteParamParserFunc[T]) *RouterParam[T] {
	return &RouterParam[T]{
		Name:   name,
		Parser: parser,
	}
}
