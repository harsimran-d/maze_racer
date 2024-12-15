package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

var playerPositions PlayerPostiions

func main() {
	playerPositions = PlayerPostiions{
		Positions:   make(map[PlayerID]Position),
		Connections: make(map[*websocket.Conn]PlayerID),
	}
	maze = NewMaze(8)
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
	go upgradehandler(conn)
}

func upgradehandler(conn *websocket.Conn) {
	pass := make(chan []byte)
	go readPump(pass, conn)
	go writePump(pass, conn)
}

func readPump(pass chan []byte, conn *websocket.Conn) {

	defer func() {
		playerId := playerPositions.Connections[conn]
		delete(playerPositions.Positions, playerId)
		delete(playerPositions.Connections, conn)
		close(pass)
	}()

	mapBytes, err := json.Marshal(maze.MazeMap)
	if err != nil {
		log.Println("error marshling game state: ", err)
		return
	}

	pass <- mapBytes
	uuid := uuid.New().String()

	playerPositions.Connections[conn] = PlayerID(uuid)
	playerPositions.Positions[PlayerID(uuid)] = Position{maze.Size / 2, maze.Size / 2}
	playerPositionsBytes, err := json.Marshal(playerPositions.Positions)

	if err != nil {
		log.Println("error marshling player positions: ", err)
		return
	}

	pass <- playerPositionsBytes

	log.Println("starting loop in read")

	for {
		messageType, moveBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("client disconnected")
				break
			}
			log.Println("error reading message: ", err)
			break
		}
		if messageType == websocket.TextMessage {
			playerId := playerPositions.Connections[conn]
			oldPosition := playerPositions.Positions[playerId]
			newPosition, isNew := maze.NewPosition(oldPosition, string(moveBytes))
			if isNew {
				log.Println(newPosition)
				playerPositions.Positions[playerId] = newPosition
				newMoveBytes, err := json.Marshal(map[string]interface{}{"playerId": playerId, "position": newPosition})
				if err != nil {
					log.Println("could not marshal new position: ", newPosition)
				} else {
					pass <- newMoveBytes
				}
			} else {

				log.Println("no new position for you: ", playerId)
			}

		} else {
			log.Printf("unsupported message type: %d", messageType)
		}
	}
}

func writePump(pass chan []byte, conn *websocket.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("error closing connection")
		}
		log.Println("connection closed")
	}()
	for messageBytes := range pass {
		err := conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			log.Println("error writing to websocket: ", err)
			break
		}
	}
}
