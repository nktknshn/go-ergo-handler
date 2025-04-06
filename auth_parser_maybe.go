package goergohandler

import (
	"context"
	"net/http"
)

// TODO: rename to tokenValidator
type tokenValidator[T any, K any] interface {
	ValidateToken(ctx context.Context, token string) (*T, bool, error)
}

type AuthParserMaybeType[T any, K any] struct {
	key         K
	tokenParser tokenParserFunc
}

type tokenParserFunc = func(ctx context.Context, r *http.Request) (string, bool, error)

func AuthParserMaybe[T any, K any](key K, tokenParser tokenParserFunc) *AuthParserMaybeType[T, K] {
	return &AuthParserMaybeType[T, K]{key, tokenParser}
}

func (a *AuthParserMaybeType[T, K]) Attach(deps tokenValidator[T, K], builder HandlerBuilder) *AttachedAuthParserMaybe[T, K] {
	attached := &AttachedAuthParserMaybe[T, K]{deps, a.tokenParser, a.key}
	builder.AddParser(attached)
	return attached
}

type AttachedAuthParserMaybe[T any, K any] struct {
	auth        tokenValidator[T, K]
	tokenParser tokenParserFunc
	key         K
}

func (a *AttachedAuthParserMaybe[T, K]) GetContextMaybe(ctx context.Context) (*T, bool) {
	data, ok := ctx.Value(a.key).(*T)
	if !ok {
		return nil, false
	}
	return data, true
}

func (a *AttachedAuthParserMaybe[T, K]) GetMaybe(r *http.Request) (*T, bool) {
	return a.GetContextMaybe(r.Context())
}

func (a *AttachedAuthParserMaybe[T, K]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	token, ok, err := a.tokenParser(ctx, r)
	if err != nil {
		return ctx, WrapError(err, defaultHttpStatusCodeErrInternal)
	}
	if !ok {
		return ctx, nil
	}
	data, ok, err := a.auth.ValidateToken(ctx, token)
	if err != nil {
		return ctx, WrapError(err, defaultHttpStatusCodeErrInternal)
	}
	if !ok {
		return ctx, nil
	}
	return context.WithValue(ctx, a.key, data), nil
}
