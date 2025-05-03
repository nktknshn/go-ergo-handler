package goergohandler

import (
	"context"
	"fmt"
	"net/http"
)

type RouterParamWithParserType[T WithParser[T]] struct {
	Name       string
	ErrMissing error
	VarsGetter VarsGetter
}

// RouterParamWithParser is same as RouterParam but it uses a parser function of the given type.
func RouterParamWithParser[T WithParser[T]](name string) *RouterParamWithParserType[T] {
	return &RouterParamWithParserType[T]{
		Name:       name,
		VarsGetter: defaultVarsGetter,
	}
}

func (rp *RouterParamWithParserType[T]) WithErrMissing(errMissing error) *RouterParamWithParserType[T] {
	rp.ErrMissing = errMissing
	return rp
}

func (rp *RouterParamWithParserType[T]) Attach(builder ParserAdder) *AttachedRouterParamWithParser[T] {
	a := &AttachedRouterParamWithParser[T]{rp}
	builder.AddParser(a)
	return a
}

type AttachedRouterParamWithParser[T WithParser[T]] struct {
	rp *RouterParamWithParserType[T]
}

func (p *AttachedRouterParamWithParser[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if p.rp.VarsGetter == nil {
		p.rp.VarsGetter = defaultVarsGetter
	}
	v, ok := p.rp.VarsGetter.GetVar(r, p.rp.Name)
	if !ok {
		err := p.rp.ErrMissing
		if err == nil {
			err = fmt.Errorf("%w: %s", ErrRouterParamMissing, p.rp.Name)
		}
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrRouterParamMissing)
	}
	var instance T
	vt, err := instance.Parse(ctx, v)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrRouterParamParsing)
	}
	err = ValidateWithValidation(vt)
	if err != nil {
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrRouterParamValidation)
	}
	return context.WithValue(ctx, routerParamKeyType(p.rp.Name), vt), nil
}

func (p *AttachedRouterParamWithParser[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedRouterParamWithParser[T]) GetContext(ctx context.Context) T {
	return GetFromContext[T](ctx, routerParamKeyType(p.rp.Name))
}
