package goergohandler

import (
	"context"
	"fmt"
	"strconv"
)

func parseBool(ctx context.Context, v string) (bool, error) {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("invalid bool value: %s", v)
	}
	return b, nil
}

func parseInt(ctx context.Context, v string) (int, error) {
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid int value: %s", v)
	}
	return i, nil
}

func parseInt64(ctx context.Context, v string) (int64, error) {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int64 value: %s", v)
	}
	return i, nil
}

func parseUInt64(ctx context.Context, v string) (uint64, error) {
	i, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid uint64 value: %s", v)
	}
	return i, nil
}

var QueryParamBool = func(name string) *QueryParamType[bool] {
	return QueryParam(name, parseBool)
}

var QueryParamBoolMaybe = func(name string) *QueryParamMaybeType[bool] {
	return QueryParamMaybe(name, parseBool)
}

var QueryParamInt = func(name string) *QueryParamType[int] {
	return QueryParam(name, parseInt)
}

var QueryParamIntMaybe = func(name string) *QueryParamMaybeType[int] {
	return QueryParamMaybe(name, parseInt)
}

var QueryParamInt64 = func(name string) *QueryParamType[int64] {
	return QueryParam(name, parseInt64)
}

var QueryParamInt64Maybe = func(name string) *QueryParamMaybeType[int64] {
	return QueryParamMaybe(name, parseInt64)
}

var QueryParamUInt64 = func(name string) *QueryParamType[uint64] {
	return QueryParam(name, parseUInt64)
}

var QueryParamUInt64Maybe = func(name string) *QueryParamMaybeType[uint64] {
	return QueryParamMaybe(name, parseUInt64)
}

var RouterParamInt64 = func(name string) *RouterParamType[int64] {
	return RouterParam(name, parseInt64)
}

var QueryParamString = func(name string) *QueryParamType[string] {
	return QueryParam(name, func(ctx context.Context, v string) (string, error) {
		return v, nil
	})
}

var QueryParamStringMaybe = func(name string) *QueryParamMaybeType[string] {
	return QueryParamMaybe(name, func(ctx context.Context, v string) (string, error) {
		return v, nil
	})
}
