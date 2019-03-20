package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) pingRouteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (s *Server) mountRoutes() {
	fs := http.FileServer(http.Dir("./webapp/build"))

	s.router.Route("/", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Handle("/", fs)
		r.Handle("/static", http.StripPrefix("/static", fs))

		r.Get("/ping", s.pingRouteHandler)

		r.Get("/transmit/{ID}", s.transmitRoute)
		r.Get("/receive/{ID}", s.receiveRoute)
	})
}
