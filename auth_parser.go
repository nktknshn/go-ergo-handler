package goergohandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	defaultHttpStatusCodeErrUnauthorized = http.StatusUnauthorized
)

var (
	// Returned when the token is missing from the request.
	ErrAuthMissingToken = errors.New("missing token")
	// Returned when token validator returned false.
	ErrAuthTokenNotFound = errors.New("token not found")
)

type tokenValidator[T any] interface {
	ValidateToken(ctx context.Context, token string) (*T, bool, error)
}

type AuthParserType[T any, K any] struct {
	key             K
	tokenParserFunc TokenParserFunc
}

type TokenParserFunc = func(ctx context.Context, r *http.Request) (string, bool, error)

// AuthParser represents a parser that parses a token string from the request.
// Attaching requires a tokenValidator that will validate the token.
// If validator returns false, ErrAuthTokenNotFound will be returned.
// If token is missing, ErrAuthMissingToken will be returned.
// If validator returns error, it will be returned wrapped with defaultHttpStatusCodeErrInternal http status code.
// On success the data returned by the validator will be set to the context with the key.
// Use WithHandlerErrorFunc to customize the error handling.
func AuthParser[T any, K any](key K, tokenParser TokenParserFunc) *AuthParserType[T, K] {
	return &AuthParserType[T, K]{key: key, tokenParserFunc: tokenParser}
}

func (a *AuthParserType[T, K]) Attach(tokenValidator tokenValidator[T], builder ParserAdder) *AttachedAuthParser[T, K] {
	attached := &AttachedAuthParser[T, K]{tokenValidator, a.tokenParserFunc, a.key}
	builder.AddParser(attached)
	return attached
}

type AttachedAuthParser[T any, K any] struct {
	tokenValidator  tokenValidator[T]
	tokenParserFunc TokenParserFunc
	key             K
}

// ParseRequest parses the request and returns the context and error.
func (a *AttachedAuthParser[T, K]) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	token, ok, err := a.tokenParserFunc(ctx, r)
	if err != nil {
		return ctx, InternalServerError(err)
	}
	if !ok {
		return ctx, WrapWithStatusCode(ErrAuthMissingToken, defaultHttpStatusCodeErrUnauthorized)
	}
	data, ok, err := a.tokenValidator.ValidateToken(ctx, token)
	if err != nil {
		return ctx, InternalServerError(err)
	}
	if !ok {
		return ctx, WrapWithStatusCode(ErrAuthTokenNotFound, defaultHttpStatusCodeErrUnauthorized)
	}
	return context.WithValue(ctx, a.key, data), nil
}

func (a *AttachedAuthParser[T, K]) GetContext(ctx context.Context) *T {
	return GetFromContext[*T](ctx, a.key)
}

func (a *AttachedAuthParser[T, K]) Get(r *http.Request) *T {
	return a.GetContext(r.Context())
}

var TokenBearerFromHeader TokenParserFunc = func(ctx context.Context, r *http.Request) (string, bool, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false, nil
	}
	// TODO: optimize
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", false, nil
	}
	token = strings.TrimPrefix(token, "bearer ")
	if token == "" {
		return "", false, nil
	}
	fmt.Println(token)
	return token, true, nil
}
