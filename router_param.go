package goergohandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	defaultHttpStatusCodeErrRouterParamParsing    = http.StatusBadRequest
	defaultHttpStatusCodeErrRouterParamMissing    = http.StatusBadRequest
	defaultHttpStatusCodeErrRouterParamValidation = http.StatusBadRequest
)

var (
	ErrRouterParamMissing = errors.New("required router param is missing")
)

func newRouterParamMissingError(paramName string) error {
	return fmt.Errorf("%w: %s", ErrRouterParamMissing, paramName)
}

type routerParamKeyType string

type RouteParamParserFunc[T any] func(ctx context.Context, v string) (T, error)

type RouterParamType[T any] struct {
	Name       string
	Parser     RouteParamParserFunc[T]
	ErrMissing error
}

func RouterParam[T any](name string, parser RouteParamParserFunc[T]) *RouterParamType[T] {
	return &RouterParamType[T]{
		Name:   name,
		Parser: parser,
	}
}

// WithMissingError sets the error to be returned if the router param is missing.
func (rp *RouterParamType[T]) WithMissingError(err error) *RouterParamType[T] {
	rp.ErrMissing = err
	return rp
}

func (rp *RouterParamType[T]) Attach(builder ParserAdder) *AttachedRouterParam[T] {
	a := &AttachedRouterParam[T]{rp}
	builder.AddParser(a)
	return a
}

type AttachedRouterParam[T any] struct {
	rp *RouterParamType[T]
}

func (p *AttachedRouterParam[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	vars := mux.Vars(r)
	v, ok := vars[p.rp.Name]
	if !ok {
		err := p.rp.ErrMissing
		if err == nil {
			err = newRouterParamMissingError(p.rp.Name)
		}
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrRouterParamMissing)
	}
	vt, err := p.rp.Parser(ctx, v)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrRouterParamParsing)
	}
	err = ValidateWithValidation(vt)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrRouterParamValidation)
	}
	return context.WithValue(ctx, routerParamKeyType(p.rp.Name), vt), nil
}

func (p *AttachedRouterParam[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedRouterParam[T]) GetContext(ctx context.Context) T {
	return GetFromContext[T](ctx, routerParamKeyType(p.rp.Name))
}
