package goergohandler

import (
	"context"
	"fmt"
)

// IgnoreContext returns a function that ignores the context argument.
// Useful for parser functions that don't need the context.
// Example:
//
// QueryParamMaybe("unpublish", geh.IgnoreContext(strconv.ParseBool))
func IgnoreContext[T any](f func(s string) (T, error)) func(ctx context.Context, s string) (T, error) {
	return func(ctx context.Context, s string) (T, error) {
		return f(s)
	}
}

// GetFromContext returns the value T stored by the key K in the context.
// If the value is not present, it panics with builderMissingKey.
// If the value is not of type T, it panics with builderCastError.
func GetFromContext[T any, K any](ctx context.Context, key K) T {
	v, ok := GetFromContextMaybe[T](ctx, key)
	if !ok {
		panic(newBuilderMissingKeyError(fmt.Sprintf("missing key from context: %v", key)))
	}
	return *v
}

// GetFromContextMaybe returns the value T stored by the key K in the context.
// If the value is not present, it returns false.
// If the value is not of type T, it panics with builderCastError.
func GetFromContextMaybe[T any, K any](ctx context.Context, key K) (*T, bool) {
	v := ctx.Value(key)
	if v == nil {
		return nil, false
	}
	casted, ok := v.(T)
	if !ok {
		panic(newBuilderCastError(fmt.Sprintf("error casting value to type %T: key: %v, value: %v, actual type: %T", *new(T), key, v, v)))
	}
	return &casted, true
}

// ValidateWithValidation validates the value v if it implements the WithValidation interface.
// If the value does not implement the WithValidation interface, it returns nil.
func ValidateWithValidation[T any](v T) error {
	validatable, ok := any(v).(WithValidation)
	if !ok {
		return nil
	}
	return validatable.Validate()
}
