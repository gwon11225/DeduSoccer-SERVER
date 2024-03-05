package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"reflect"
)

var (
	rc = RoomClient{make(map[string]*Room)}
)

type RoomClient struct {
	roomMap map[string]*Room
}

func (rc RoomClient) Enter(roomId string, username string, conn *websocket.Conn) {
	room, exist := rc.roomMap[roomId]
	if exist {
		if reflect.DeepEqual(room.rUser, User{}) {
			rc.roomMap[roomId].SetUser(User{conn, username, -460, -224}, "RED")
			rc.sendMessage(room.bUser.conn, fmt.Sprintf("ENTER/%s/%s", username, "RED"))
			rc.sendMessage(room.bUser.conn, fmt.Sprintf("ENTER/%s/%s", room.bUser.name, "BLUE"))
		} else {
			rc.roomMap[roomId].SetUser(User{conn, username, 460, -224}, "BLUE")
			rc.sendMessage(room.rUser.conn, fmt.Sprintf("ENTER/%s/%s", username, "BLUE"))
			rc.sendMessage(room.bUser.conn, fmt.Sprintf("ENTER/%s/%s", room.rUser.name, "RED"))
		}
	} else {
		rc.roomMap[roomId] = &Room{bScore: 0, rScore: 0}
		rc.roomMap[roomId].SetUser(User{conn, username, -460, -224}, "RED")
		rc.roomMap[roomId].SetBall(Ball{0, 80})
	}
}

func (rc RoomClient) Quit(roomId string, username string) {
	room := rc.roomMap[roomId]

	if UserCheck(*room, username) == "RED" {
		rc.roomMap[roomId].DelUser("RED")
		if !reflect.DeepEqual(rc.roomMap[roomId].bUser, User{}) {
			rc.roomMap[roomId].DelUser("BLUE")
		}
	} else {
		rc.roomMap[roomId].DelUser("BLUE")
		if !reflect.DeepEqual(rc.roomMap[roomId].bUser, User{}) {
			rc.roomMap[roomId].DelUser("RED")
		}
	}

	room = rc.roomMap[roomId]

	if reflect.DeepEqual(room.rUser, User{}) && reflect.DeepEqual(room.bUser, User{}) {
		delete(rc.roomMap, roomId)
	}
}

func (rc RoomClient) Goal(roomId string, team string) {
	room := rc.roomMap[roomId]

	if team == "RED" {
		rc.roomMap[roomId].Goal("RED")
		rc.sendMessage(room.bUser.conn, fmt.Sprintf("GOAL/RED/%d", rc.roomMap[roomId].rScore))
	} else {
		rc.roomMap[roomId].Goal("BLUE")
		rc.sendMessage(room.rUser.conn, fmt.Sprintf("GOAL/BLUE/%d", rc.roomMap[roomId].bScore))
	}
}

func (rc RoomClient) Out(roomId string, team string) {
	room := rc.roomMap[roomId]
	if team == "RED" {
		rc.sendMessage(room.bUser.conn, fmt.Sprintf("OUT/%s", team))
	} else {
		rc.sendMessage(room.rUser.conn, fmt.Sprintf("OUT/%s", team))
	}
}

func (rc RoomClient) Move(roomId string, username string, posX int, posY int) {
	if username == "BALL" {
		rc.roomMap[roomId].MoveBall(posX, posY)
	} else {
		rc.roomMap[roomId].MoveUser(username, posX, posY)
	}
	room := rc.roomMap[roomId]
	if UserCheck(*rc.roomMap[roomId], username) == "RED" {
		rc.sendMessage(room.bUser.conn, fmt.Sprintf("MOVE/%s/%d/%d", room.rUser.name, room.rUser.posX, room.rUser.posY))
	} else {
		rc.sendMessage(room.rUser.conn, fmt.Sprintf("MOVE/%s/%d/%d", room.bUser.name, room.bUser.posX, room.bUser.posY))
	}
}

func (rc RoomClient) sendMessage(conn *websocket.Conn, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		conn.Close()
	}
}
