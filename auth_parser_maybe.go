package goergohandler

import (
	"context"
	"net/http"
)

type AuthParserMaybeType[T any, K any] struct {
	key             K
	tokenParserFunc TokenParserFunc
}

// AuthParserMaybe is the same as AuthParser but it allows the token to be missing or
// validator to return false.
func AuthParserMaybe[T any, K any](key K, tokenParser TokenParserFunc) *AuthParserMaybeType[T, K] {
	return &AuthParserMaybeType[T, K]{key, tokenParser}
}

func (a *AuthParserMaybeType[T, K]) Attach(tokenValidator tokenValidator[T], builder ParserAdder) *AttachedAuthParserMaybe[T, K] {
	attached := &AttachedAuthParserMaybe[T, K]{tokenValidator, a.tokenParserFunc, a.key}
	builder.AddParser(attached)
	return attached
}

type AttachedAuthParserMaybe[T any, K any] struct {
	auth            tokenValidator[T]
	tokenParserFunc TokenParserFunc
	key             K
}

func (a *AttachedAuthParserMaybe[T, K]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	token, ok, err := a.tokenParserFunc(ctx, r)
	if err != nil {
		return ctx, InternalServerError(err)
	}
	if !ok {
		return ctx, nil
	}
	data, ok, err := a.auth.ValidateToken(ctx, token)
	if err != nil {
		return ctx, InternalServerError(err)
	}
	if !ok {
		return ctx, nil
	}
	return context.WithValue(ctx, a.key, data), nil
}

func (a *AttachedAuthParserMaybe[T, K]) GetContextMaybe(ctx context.Context) (*T, bool) {
	return GetFromContextMaybe[*T](ctx, a.key)
}

func (a *AttachedAuthParserMaybe[T, K]) GetMaybe(r *http.Request) (*T, bool) {
	return a.GetContextMaybe(r.Context())
}
