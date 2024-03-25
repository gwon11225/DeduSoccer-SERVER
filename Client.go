package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"reflect"
	"sync"
	"time"
)

var (
	rc    = RoomClient{make(map[string]*Room)}
	mutex = sync.Mutex{}
)

type RoomClient struct {
	roomMap map[string]*Room
}

func (rc RoomClient) Enter(roomId string, username string, conn *websocket.Conn) {
	room, exist := rc.roomMap[roomId]
	if exist {
		if reflect.DeepEqual(room.rUser, User{}) {
			rc.roomMap[roomId].SetUser(User{conn, username, -4.6, -2.24}, "RED")
			log.Println("ENTER RED USER")
			rc.sendMessage(room.rUser.conn, "RED")
			rc.sendMessage(room.bUser.conn, fmt.Sprintf("ENTER/%s/%s", username, "RED"))
			rc.sendMessage(room.rUser.conn, fmt.Sprintf("ENTER/%s/%s", room.bUser.name, "BLUE"))
		} else {
			log.Println("ENTER BLUE USER")
			rc.roomMap[roomId].SetUser(User{conn, username, 4.6, 0}, "BLUE")
			rc.sendMessage(room.bUser.conn, "BLUE")
			rc.sendMessage(room.rUser.conn, fmt.Sprintf("ENTER/%s/%s", username, "BLUE"))
			rc.sendMessage(room.bUser.conn, fmt.Sprintf("ENTER/%s/%s", room.rUser.name, "RED"))
		}
	} else {
		log.Println("ENTER RED USER")
		rc.roomMap[roomId] = &Room{bScore: 0, rScore: 0}
		rc.roomMap[roomId].SetUser(User{conn, username, -4.6, 0}, "RED")
		rc.roomMap[roomId].SetBall(Ball{0.0, 1.0, 0.0, 0.0})
		rc.sendMessage(rc.roomMap[roomId].rUser.conn, "RED")
	}
}

func (rc RoomClient) Quit(roomId string, username string) {
	room := rc.roomMap[roomId]
	user := UserCheck(*room, username)

	if user == "RED" {
		rc.roomMap[roomId].DelUser("RED")
		if !reflect.DeepEqual(rc.roomMap[roomId].bUser, User{}) {
			rc.sendMessage(room.bUser.conn, fmt.Sprintf("QUIT/%s", username))
		}
	} else if user == "BLUE" {
		rc.roomMap[roomId].DelUser("BLUE")
		if !reflect.DeepEqual(rc.roomMap[roomId].rUser, User{}) {
			rc.sendMessage(room.rUser.conn, fmt.Sprintf("QUIT/%s", username))
		}
	}

	log.Println("USER QUIT username : " + username)

	room = rc.roomMap[roomId]

	if reflect.DeepEqual(room.rUser, User{}) && reflect.DeepEqual(room.bUser, User{}) {
		delete(rc.roomMap, roomId)
		log.Println("ROOM DELETE roomId : " + roomId)
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

func (rc RoomClient) Move(roomId string, username string, posX float64, posY float64) {
	rc.roomMap[roomId].MoveUser(username, posX, posY)

	room := rc.roomMap[roomId]
	if UserCheck(*rc.roomMap[roomId], username) == "RED" {
		if !reflect.DeepEqual(room.bUser, User{}) {
			rc.sendMessage(room.bUser.conn, fmt.Sprintf("MOVE/%s/%f/%f", room.rUser.name, room.rUser.posX, room.rUser.posY))
		}
	} else {
		if !reflect.DeepEqual(room.rUser, User{}) {
			rc.sendMessage(room.rUser.conn, fmt.Sprintf("MOVE/%s/%f/%f", room.bUser.name, room.bUser.posX, room.bUser.posY))
		}
	}
}

func (rc RoomClient) Ball(roomId string) {
	ticker := time.Tick(time.Second / 60)
	for {
		for range ticker {
			room, exist := rc.roomMap[roomId]

			if !exist {
				return
			}

			room.BallUpdate()
			room.CollisionFloor()

			if !reflect.DeepEqual(room.rUser, User{}) {
				go rc.sendMessage(room.rUser.conn, fmt.Sprintf("MOVE/BALL/%f/%f", room.ball.posX, room.ball.posY))
			}

			if !reflect.DeepEqual(room.bUser, User{}) {
				go rc.sendMessage(room.bUser.conn, fmt.Sprintf("MOVE/BALL/%f/%f", room.ball.posX, room.ball.posY))
			}
		}
	}
}

func (rc RoomClient) Coll(roomId string, posX float64, posY float64) {
	rc.roomMap[roomId].CollisionUser(posX, posY)
}

func (rc RoomClient) sendMessage(conn *websocket.Conn, message string) {
	mutex.Lock()
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		conn.Close()
	}
	mutex.Unlock()
}
