package main

import "github.com/gorilla/websocket"

var (
	rm = make(RoomMap)
)

const (
	BALL EntityType = iota
	USER
)

const (
	RED Team = iota
	BLUE
)

type EntityType int
type Team int
type RoomMap map[string]Room

type Entity struct {
	eType EntityType
	posX  string
	posY  string
}

type Room struct {
	roomId string
	user   *[]User
	ball   *Entity
	rScore *int
	bScore *int
}

type User struct {
	Entity
	conn     *websocket.Conn
	username string
	team     Team
}

func (e *Entity) move(posX string, posY string) {
	e.posX = posX
	e.posY = posY
}

func (r Room) AddUser(user User) {
	*r.user = append(*r.user, user)
}

func (r Room) GetUser(username string) *User {
	for _, user := range *r.user {
		if user.username == username {
			return &user
		}
	}
	return &User{}
}

func (r Room) DeleteUser(username string) {
	for index, user := range *r.user {
		if user.username == username {
			*r.user = append((*r.user)[:index], (*r.user)[index+1:]...)
			return
		}
	}
}

func (r Room) RedGoal() {
	*r.rScore += 1
}

func (r Room) BlueGoal() {
	*r.bScore += 1
}
