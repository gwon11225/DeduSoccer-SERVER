package main

import (
	"github.com/gorilla/websocket"
	"reflect"
)

type Room struct {
	rUser  User
	bUser  User
	ball   Ball
	bScore int
	rScore int
}

type User struct {
	conn *websocket.Conn
	name string
	posX int
	posY int
}

type Ball struct {
	posX int
	posY int
}

func (r *Room) SetUser(user User, team string) {
	if team == "RED" {
		r.rUser = user
	} else {
		r.bUser = user
	}
}

func (r *Room) DelUser(team string) {
	if team == "RED" {
		r.rUser = User{}
	} else {
		r.bUser = User{}
	}
}

func (r *Room) SetBall(ball Ball) {
	r.ball = ball
}

func (r *Room) Goal(team string) {
	if team == "RED" {
		r.rScore++
	} else {
		r.bScore++
	}
}

func (r *Room) MoveUser(username string, posX int, posY int) {
	team := UserCheck(*r, username)
	if team == "RED" {
		r.rUser.posX = posX
		r.rUser.posY = posY
	} else {
		r.bUser.posX = posX
		r.bUser.posY = posY
	}
}

func (r *Room) MoveBall(posX int, posY int) {
	r.ball.posX = posX
	r.ball.posY = posY
}

func UserCheck(r Room, username string) string {
	if !reflect.DeepEqual(r.rUser, User{}) && r.rUser.name == username {
		return "RED"
	}
	if !reflect.DeepEqual(r.bUser, User{}) && r.bUser.name == username {
		return "BLUE"
	}

	return "NONE"
}
