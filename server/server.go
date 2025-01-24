package server

import (
	"github.com/flint92/webapp/context"
	"github.com/flint92/webapp/handler"
	"log"
	"net/http"
)

type Server interface {
	handler.Routable
	Name() string
	Start(address string)
}

func NewServer(name string) Server {
	return &sdkHttpServer{
		name:    name,
		handler: handler.NewHandler(),
	}
}

type sdkHttpServer struct {
	name    string
	handler handler.Handler
}

func (s *sdkHttpServer) Name() string {
	return s.name
}

func (s *sdkHttpServer) Route(method string, pattern string, handler func(*context.Context)) {
	s.handler.Route(method, pattern, handler)
}

func (s *sdkHttpServer) Start(address string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(w, r)
		s.handler.ServeHTTP(ctx)
	})
	log.Printf("Starting server=[%s] at %s\n", s.Name(), address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
