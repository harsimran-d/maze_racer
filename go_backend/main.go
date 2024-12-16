package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Position struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}
type PlayerID string

type PlayerPostiions struct {
	Positions   map[PlayerID]Position
	Connections map[*websocket.Conn]PlayerID
}

var maze *Maze
var hub *Hub
var playerPositions PlayerPostiions

func main() {
	playerPositions = PlayerPostiions{
		Positions:   make(map[PlayerID]Position),
		Connections: make(map[*websocket.Conn]PlayerID),
	}
	maze = NewMaze(25)
	hub = NewHub()
	go hub.Run()
	router := gin.New()
	router.GET("/", roothandler)
	router.GET("/game", gameHandler)
	router.Run(":3000")
}

func roothandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello world"})
}

func gameHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("could not upgrade connection: ", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	hub.register <- conn
}
