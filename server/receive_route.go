package server

import (
	"net/http"

	"github.com/go-chi/chi"

	"bitbucket.org/emotech/common/golang/logs"

	"github.com/gorilla/websocket"
)

func (s *Server) receiveRoute(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "ID")
	logs.Printf("GET /receive request with ID [%s]", ID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Printf("error upgrading receive request to websocket: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for {
		ch := s.DataHandler.ReceiveData(ID)
		if ch == nil {
			logs.Printf("no transmission channel found for ID [%s]; closing connection", ID)
			if err := conn.Close(); err != nil {
				logs.Printf("error closing websocket connection: %v", err)
				return
			}

		}
		data, ok := <-ch
		if !ok {
			logs.Printf("transmission channel is closed form the other side, closing ws ...")
			if err := conn.Close(); err != nil {
				logs.Printf("error closing connection: %v", err)
			}
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			logs.Printf("error writing to channel: %v for id [%s]", err, ID)
			break
		}
	}

}
