package server

import (
	"github.com/flint92/webapp/context"
	"github.com/flint92/webapp/filter"
	"github.com/flint92/webapp/handler"
	"log"
	"net/http"
)

type Server interface {
	handler.Routable
	Name() string
	Start(address string)
}

func NewServer(name string, builders ...filter.FilterBuilder) Server {
	h := handler.NewHandler()

	var root filter.Filter = func(ctx *context.Context) {
		h.ServeHTTP(ctx)
	}

	for i := len(builders) - 1; i >= 0; i-- {
		root = builders[i](root)
	}

	return &sdkHttpServer{
		name:    name,
		handler: h,
		root:    root,
	}
}

type sdkHttpServer struct {
	name    string
	handler handler.Handler
	root    filter.Filter
}

func (s *sdkHttpServer) Name() string {
	return s.name
}

func (s *sdkHttpServer) Route(method string, pattern string, handler handler.HandlerFunc) {
	s.handler.Route(method, pattern, handler)
}

func (s *sdkHttpServer) Start(address string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(w, r)
		s.root(ctx)
	})
	log.Printf("Starting server=[%s] at %s\n", s.Name(), address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
