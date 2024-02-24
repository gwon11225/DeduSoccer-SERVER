package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

var (
	upgrader = websocket.Upgrader{}
)

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	username := r.URL.Query().Get("username")
	conn, err := upgrader.Upgrade(w, r, nil)

	if roomId == "" || username == "" {
		http.Error(w, "Param is nil", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, exists := rm[roomId]

	user := User{Entity: Entity{eType: USER, posX: "0", posY: "0"}, username: username, conn: conn, team: RED}

	conn.SetCloseHandler(func(code int, text string) error {
		sendMessage(roomId, fmt.Sprintf("QUIT/%s", username))
		rm[roomId].DeleteUser(user.username)
		log.Println(fmt.Sprintf("QUIT ROOM. ROOMID : %s, WHO : %s", roomId, username))
		if len(*rm[roomId].user) == 0 {
			delete(rm, roomId)
			log.Println(fmt.Sprintf("DEL ROOM. ROOMID : %s", roomId))
		}
		return nil
	})

	if !exists {
		users := make([]User, 0)
		ball := Entity{eType: BALL, posX: "0", posY: "0"}
		rScore := 0
		bScore := 0
		rm[roomId] = Room{roomId: roomId, ball: &ball, user: &users, rScore: &rScore, bScore: &bScore}
		rm[roomId].AddUser(user)
		log.Println(fmt.Sprintf("MAKE ROOM. ROOMID : %s", roomId))
	} else {
		if len(*rm[roomId].user)%2 == 1 {
			user.team = BLUE
		}

		rm[roomId].AddUser(user)
	}

	go messageBroker(conn, roomId, username, user.team)
}

func messageBroker(conn *websocket.Conn, roomId string, username string, team Team) {
	uTeam := "RED"
	if team == BLUE {
		uTeam = "BLUE"
	}
	sendMessage(roomId, fmt.Sprintf("ENTER/%s/%s", username, uTeam))
	log.Println(fmt.Sprintf("ENTER ROOM. ROOMID : %s, WHO : %s", roomId, username))
	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			conn.Close()
			return
		}

		sMessage := strings.Split(string(message), "/")
		if sMessage[0] == "MOVE" {
			if sMessage[1] == "BALL" {
				rm[roomId].ball.move(sMessage[2], sMessage[3])
			} else {
				(*rm[roomId].GetUser(sMessage[1])).move(sMessage[2], sMessage[3])
			}
			sendMessage(roomId, fmt.Sprintf("%s/%s/%s/%s", "MOVE", sMessage[1], sMessage[2], sMessage[3]))
		} else if sMessage[0] == "GOAL" {
			room := rm[roomId]
			score := 0
			log.Println("11111111")
			if sMessage[2] == "RED" {
				log.Println("11111111")
				rm[roomId].RedGoal()
				log.Println("11111111")
				log.Println(*room.rScore)
				score = *room.rScore
			} else {
				rm[roomId].BlueGoal()
				log.Println(*room.bScore)
				score = *room.bScore
			}
			log.Println("122222222")
			sendMessage(roomId, fmt.Sprintf("GOAL/%s/%s/%d", sMessage[1], sMessage[2], score))
		} else if sMessage[0] == "OUT" {
			sendMessage(roomId, fmt.Sprintf("OUT/%s", sMessage[1]))
		} else {
			sendMessage(roomId, "잘못 보냈다")
		}
	}
}

func sendMessage(roomId string, message string) {
	for _, user := range *rm[roomId].user {
		err := user.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			user.conn.Close()
			return
		}
	}
}
