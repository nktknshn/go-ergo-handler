package goergohandler

import (
	"context"
	"net/http"
)

type CorsParser struct {
	allowedOrigins []string
}

func NewCorsParser(allowedOrigins []string) *CorsParser {
	return &CorsParser{allowedOrigins}
}

func (p *CorsParser) Attach(builder HandlerBuilder) *AttachedCorsParser {
	attached := &AttachedCorsParser{p}
	builder.AddParser(attached)
	return attached
}

type AttachedCorsParser struct {
	corsParser *CorsParser
}

func (p *AttachedCorsParser) ParseRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	// TODO
	// origin := r.Header.Get("Origin")
	// if origin == "" {
	// 	return ctx, nil
	// }
	// for _, allowedOrigin := range p.corsParser.allowedOrigins {
	// 	if allowedOrigin == origin {
	// 		w.Header().Set("Access-Control-Allow-Origin", origin)
	// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 		return ctx, nil
	// 	}
	// }
	return ctx, nil
}
