package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"strings"
)

var (
	upgrader = websocket.Upgrader{}
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	username := r.URL.Query().Get("username")

	if roomId == "" || username == "" {
		http.Error(w, "Param is nil", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		w.Write([]byte("ㅋㅋ 이게 맞는거 같아??"))
	}

	conn.SetCloseHandler(func(code int, text string) error {
		rc.Quit(roomId, username)
		return nil
	})

	rc.Enter(roomId, username, conn)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			conn.Close()
			return
		}

		sMessage := strings.Split(string(message), "/")

		switch sMessage[0] {
		case "GOAL":
			rc.Goal(roomId, sMessage[1])
		case "OUT":
			rc.Out(roomId, sMessage[1])
		case "MOVE":
			posX, _ := strconv.Atoi(sMessage[2])
			posY, _ := strconv.Atoi(sMessage[3])
			rc.Move(roomId, sMessage[1], posX, posY)
		case "QUIT":
			conn.Close()
			return
		}
	}
}
