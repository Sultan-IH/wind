package server

import (
	"fmt"
	"net/http"

	logs "log"

	"github.com/Sultan-IH/wind/datahandler"
	"github.com/Sultan-IH/wind/plug"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

type Server struct {
	Config      *Config
	router      chi.Router
	DataHandler *datahandler.DataHandler
}

func NewServer() (*Server, error) {
	c := &Config{}
	if err := c.Parse(); err != nil {
		return nil, err
	}
	err := plug.GetToken(c.TPLinkUsername, c.TPLinkPwd, c.TPLinkUUID)
	if err != nil {
		return nil, nil
	}
	plugs, err := plug.GetDeviceList()
	if err != nil {
		return nil, nil
	}
	dh := datahandler.NewDataHandler(plugs)

	s := &Server{
		Config:      c,
		router:      chi.NewRouter(),
		DataHandler: dh,
	}
	s.mountRoutes()

	return s, nil
}

func (s *Server) Run() {
	logs.Println("WIND started and now running on port:", s.Config.Port)
	logs.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), s.router))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  (2 ^ 10) * 4,
	WriteBufferSize: (2 ^ 10) * 4,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
