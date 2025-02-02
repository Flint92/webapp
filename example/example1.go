package main

import (
	"github.com/flint92/webapp/context"
	"github.com/flint92/webapp/filter"
	"github.com/flint92/webapp/server"
	"net/http"
)

func helloWorldHandler(ctx *context.Context) {
	_ = ctx.OkJson("hello World")
}

func helloHandler(ctx *context.Context) {
	_ = ctx.OkJson("hello " + ctx.PathParams["username"])
}

func helloAnyHandler(ctx *context.Context) {
	_ = ctx.OkJson("hello any")
}

func main() {
	s1 := server.NewServer("test-server", filter.MetricFilterBuilder)
	s1.Route(http.MethodGet, "/hello", helloWorldHandler)
	s1.Route(http.MethodGet, "/hello/:username", helloHandler)
	s1.Route(http.MethodGet, "/hello/v1/*", helloAnyHandler)
	s1.Start(":8080")
}
