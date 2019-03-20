package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"

	"bitbucket.org/emotech/common/golang/logs"
)

const (
	pongWait   = time.Second * 2
	pingPeriod = pongWait / 2
)

func (s *Server) transmitRoute(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "ID")
	logs.Printf("GET /transmit with ID [%v]", ID)

	if ok := s.DataHandler.StartTransmission(ID); !ok {
		logs.Printf("error starting transmission for ID [%s]", ID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	defer func() {
		if err := conn.Close(); err != nil {
			logs.Printf("error closing websocket connection: %v", err)
		}
		s.DataHandler.EndTransmission(ID)
	}()

	if err != nil {
		logs.Printf("error upgrading transmit request to websocket: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPingHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		logs.Printf("PING hanlder here")
		return nil
	})
	for {
		logs.Printf("reading msgs from socket ...")
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logs.Printf("error reading msg from transmit socket: %v", err)
			break
		}

		logs.Printf("received msg [%s]", string(msg))
		if msgType != websocket.TextMessage {
			logs.Print("error reading msg from transmit socket invalid msg type")
			break
		}
		s.DataHandler.RecordData(ID, msg)
	}

}
