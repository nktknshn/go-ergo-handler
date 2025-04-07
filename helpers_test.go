package goergohandler_test

import (
	"context"
	"errors"
	"testing"

	goergohandler "github.com/nktknshn/go-ergo-handler"
	"github.com/stretchr/testify/require"
)

type keyType string

var key keyType = "context_key"

func TestGetFromContext(t *testing.T) {
	ctxWithoutValue := context.Background()
	ctxWithValue := context.WithValue(ctxWithoutValue, key, 10)
	ctxWithWrongValue := context.WithValue(ctxWithoutValue, key, "10")

	require.Equal(t, 10, goergohandler.GetFromContext[int](ctxWithValue, key))

	require.Panics(t, func() {
		goergohandler.GetFromContext[int](ctxWithoutValue, key)
	})

	require.Panics(t, func() {
		goergohandler.GetFromContext[int](ctxWithWrongValue, key)
	})
}

func TestGetFromContextMaybe(t *testing.T) {
	ctxWithoutValue := context.Background()
	ctxWithValue := context.WithValue(ctxWithoutValue, key, 10)
	ctxWithWrongValue := context.WithValue(ctxWithoutValue, key, "10")

	require.Equal(t, 10, goergohandler.GetFromContext[int](ctxWithValue, key))
	v, ok := goergohandler.GetFromContextMaybe[int](ctxWithoutValue, key)
	require.Equal(t, 0, v)
	require.False(t, ok)

	require.Panics(t, func() {
		goergohandler.GetFromContextMaybe[int](ctxWithWrongValue, key)
	})
}

type testStruct1 struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func (t testStruct1) Validate() error {
	if t.Field1 == "" {
		return errors.New("field1 is required")
	}
	return nil
}

type testStruct2 struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestValidateWithValidation(t *testing.T) {
	require.NoError(t, goergohandler.ValidateWithValidation(testStruct1{Field1: "test"}))
	require.Error(t, goergohandler.ValidateWithValidation(testStruct1{}))
	require.NoError(t, goergohandler.ValidateWithValidation(testStruct2{}))
}
