package handler

import (
	"github.com/flint92/webapp/context"
	"log"
)

type Routable interface {
	Route(method string, pattern string, handler func(*context.Context))
}

type Handler interface {
	Routable
	ServeHTTP(ctx *context.Context)
}

type BasedOnMapHandler struct {
	routes map[string]func(*context.Context)
}

func NewHandler() Handler {
	return &BasedOnMapHandler{
		routes: make(map[string]func(*context.Context)),
	}
}

func (h *BasedOnMapHandler) Route(method string, pattern string, handler func(*context.Context)) {
	k := routeKey(method, pattern)
	if _, ok := h.routes[k]; ok {
		log.Fatalf("route already exists: %s", k)
	} else {
		h.routes[k] = handler
	}
}

func (h *BasedOnMapHandler) ServeHTTP(ctx *context.Context) {
	k := routeKey(ctx.R.Method, ctx.R.URL.Path)
	if handler, ok := h.routes[k]; ok {
		handler(ctx)
	} else {
		err := ctx.NotFoundJson()
		if err != nil {
			log.Printf("error writing json: %v", err)
		}
	}
}

func routeKey(method, path string) string {
	return method + "#" + path
}
