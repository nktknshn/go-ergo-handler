package goergohandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	defaultHttpStatusCodeErrQueryParamMissing    = http.StatusBadRequest
	defaultHttpStatusCodeErrQueryParamParsing    = http.StatusBadRequest
	defaultHttpStatusCodeErrQueryParamValidation = http.StatusBadRequest
)

var (
	ErrQueryParamMissing = errors.New("required query param is missing")
)

func newQueryParamMissingError(paramName string) error {
	return fmt.Errorf("%w: %s", ErrQueryParamMissing, paramName)
}

type queryParamKeyType string

type QueryParamParserFunc[T any] func(ctx context.Context, v string) (T, error)

type QueryParamType[T any] struct {
	Name       string
	Parser     QueryParamParserFunc[T]
	ErrMissing error
}

// QueryParam is a parser that parses a required query param from the request.
// If the query param is missing, it returns ErrQueryParamMissing.
// If the type implements WithValidation, it will be validated.
func QueryParam[T any](
	name string,
	parser QueryParamParserFunc[T],
) *QueryParamType[T] {
	return &QueryParamType[T]{
		Name:   name,
		Parser: parser,
	}
}

// WithMissingError sets the error to be returned if the query param is missing.
func (qp *QueryParamType[T]) WithMissingError(err error) *QueryParamType[T] {
	qp.ErrMissing = err
	return qp
}

func (qp *QueryParamType[T]) Attach(b ParserAdder) *AttachedQueryParam[T] {
	a := &AttachedQueryParam[T]{qp}
	b.AddParser(a)
	return a
}

type AttachedQueryParam[T any] struct {
	qp *QueryParamType[T]
}

func (p *AttachedQueryParam[T]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if !r.URL.Query().Has(p.qp.Name) {
		err := p.qp.ErrMissing
		if err == nil {
			err = newQueryParamMissingError(p.qp.Name)
		}
		return ctx, WrapWithStatusCode(err, defaultHttpStatusCodeErrQueryParamMissing)
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
	return context.WithValue(ctx, queryParamKeyType(p.qp.Name), v), nil
}

func (p *AttachedQueryParam[T]) Get(r *http.Request) T {
	return p.GetContext(r.Context())
}

func (p *AttachedQueryParam[T]) GetContext(ctx context.Context) T {
	return GetFromContext[T](ctx, queryParamKeyType(p.qp.Name))
}
