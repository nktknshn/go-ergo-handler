package goergohandler

import (
	"context"
	"net/http"
)

type AttachabaleMiddlewareType struct {
	MiddlewareFunc MiddlewareFunc
}

func AttachabaleMiddleware(m MiddlewareFunc) *AttachabaleMiddlewareType {
	return &AttachabaleMiddlewareType{m}
}

func (m *AttachabaleMiddlewareType) Attach(b ParserAdder) *AttachedAttachabaleMiddlewareType {
	am := &AttachedAttachabaleMiddlewareType{m}
	b.AddParser(am)
	return am
}

type AttachedAttachabaleMiddlewareType struct {
	m *AttachabaleMiddlewareType
}

func (m *AttachedAttachabaleMiddlewareType) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	var nextContext context.Context
	var hextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if not called we stop propagation
		// otherwrise we extract context from the request and return it!!!!!
		nextContext = r.Context()
	})
	m.m.MiddlewareFunc(hextHandler).ServeHTTP(w, r)
	if nextContext == nil {
		return ctx, ErrStopPropagation
	}
	return nextContext, nil
}
