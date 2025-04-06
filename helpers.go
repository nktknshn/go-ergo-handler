package goergohandler

import "context"

func IgnoreContext[T any](f func(s string) (T, error)) func(ctx context.Context, s string) (T, error) {
	return func(ctx context.Context, s string) (T, error) {
		return f(s)
	}
}
