package handler

import (
	"github.com/flint92/webapp/context"
)

type HandlerFunc func(ctx *context.Context)

type Routable interface {
	Route(method string, pattern string, handler HandlerFunc)
}

type Handler interface {
	Routable
	ServeHTTP(ctx *context.Context)
}

func NewHandler() Handler {
	return &BasedOnTreeHandler{
		getRoot:     newRootNode("/"),
		postRoot:    newRootNode("/"),
		putRoot:     newRootNode("/"),
		patchRoot:   newRootNode("/"),
		deleteRoot:  newRootNode("/"),
		optionsRoot: newRootNode("/"),
		headRoot:    newRootNode("/"),
		connectRoot: newRootNode("/"),
		traceRoot:   newRootNode("/"),
	}
}
