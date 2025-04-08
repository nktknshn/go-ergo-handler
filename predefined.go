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
