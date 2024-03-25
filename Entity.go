package main

import (
	"github.com/gorilla/websocket"
	"log"
	"reflect"
	"time"
)

const gravity float64 = 0.5
const restitutionY float64 = 0.99
const restitutionX float64 = 0.995

type Room struct {
	rUser          User
	bUser          User
	ball           Ball
	bScore         int
	rScore         int
	lastUpdateTime time.Time
}

type User struct {
	conn *websocket.Conn
	name string
	posX float64
	posY float64
}

type Ball struct {
	posX      float64
	posY      float64
	velocityY float64
	velocityX float64
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

func (r *Room) MoveUser(username string, posX float64, posY float64) {
	team := UserCheck(*r, username)
	if team == "RED" {
		r.rUser.posX = posX
		r.rUser.posY = posY
	} else if team == "BLUE" {
		r.bUser.posX = posX
		r.bUser.posY = posY
	}
}

func (r *Room) MoveBall(posX float64, posY float64) {
	r.ball.posX = posX
	r.ball.posY = posY
}

func (r *Room) BallUpdate() {
	deltaTime := (time.Second / 60).Seconds()
	r.ball.velocityY += gravity * deltaTime
	r.ball.posX += r.ball.velocityX * deltaTime

	if r.ball.velocityY >= r.ball.posY {
		r.ball.posY = 0
	} else {
		r.ball.posY -= r.ball.velocityY
	}

	r.ball.velocityX *= restitutionX
	if r.ball.velocityX < 0.1 {
		r.ball.velocityX = 0
	}
}

func (r *Room) CollisionUser(posX float64, posY float64) {
	log.Println(posX, posY)
	if posX <= -0.3 {
		r.ball.velocityX = -10.0
	} else if posX >= 0.3 {
		r.ball.velocityX = 10.0
	}

	if posY <= -0.3 {
		r.ball.velocityY = 0.1
	} else if posY >= 0.3 {
		r.ball.velocityY = -0.1
	}
}

func (r *Room) CollisionFloor() {
	if r.ball.posY <= 0 {
		r.ball.velocityY = -r.ball.velocityY * restitutionY
	}
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
