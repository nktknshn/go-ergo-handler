package goergohandler

import (
	"context"
	"errors"
	"net/http"
)

const (
	defaultHttpStatusCodeErrUnauthorized = http.StatusUnauthorized
	defaultHttpStatusCodeErrInternal     = http.StatusInternalServerError
)

var (
	ErrNoToken = errors.New("no token")
	ErrNoUser  = errors.New("no user")
)

type AuthParserType[T any, K any] struct {
	key         K
	tokenParser tokenParserFunc
}

func AuthParser[T any, K any](key K, tokenParser tokenParserFunc) *AuthParserType[T, K] {
	return &AuthParserType[T, K]{key, tokenParser}
}

func (a *AuthParserType[T, K]) Attach(deps tokenValidator[T, K], builder HandlerBuilder) *AttachedAuthParser[T, K] {
	attached := &AttachedAuthParser[T, K]{deps, a.tokenParser, a.key}
	builder.AddParser(attached)
	return attached
}

type AttachedAuthParser[T any, K any] struct {
	auth        tokenValidator[T, K]
	tokenParser tokenParserFunc
	key         K
}

func (a *AttachedAuthParser[T, K]) GetUser(ctx context.Context) *T {
	data, ok := ctx.Value(a.key).(*T)
	if !ok {
		return nil
	}
	return data
}

func (a *AttachedAuthParser[T, K]) GetUserRequest(r *http.Request) *T {
	return a.GetUser(r.Context())
}

func (a *AttachedAuthParser[T, K]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	token, ok, err := a.tokenParser(ctx, r)
	if err != nil {
		return ctx, WrapError(err, defaultHttpStatusCodeErrInternal)
	}
	if !ok {
		return ctx, WrapError(ErrNoToken, defaultHttpStatusCodeErrUnauthorized)
	}
	data, ok, err := a.auth.ValidateToken(ctx, token)
	if err != nil {
		return ctx, WrapError(err, defaultHttpStatusCodeErrInternal)
	}
	if !ok {
		return ctx, WrapError(ErrNoUser, defaultHttpStatusCodeErrUnauthorized)
	}
	return context.WithValue(ctx, a.key, data), nil
}
