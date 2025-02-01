package main

import (
	"github.com/flint92/webapp/context"
	"github.com/flint92/webapp/filter"
	"github.com/flint92/webapp/server"
	"net/http"
)

func helloHandler(ctx *context.Context) {
	_ = ctx.OkJson("hello world")
}

func main() {
	s1 := server.NewServer("test-server", filter.MetricFilterBuilder)
	s1.Route(http.MethodGet, "/hello", helloHandler)
	s1.Start(":8080")
}
