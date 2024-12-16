package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	broadcast  chan []byte
	conns      map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 256),
		conns:      make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn, 20),
		unregister: make(chan *websocket.Conn, 20),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			log.Println("a player regiestered")
			mapBytes, err := json.Marshal(maze)
			if err != nil {
				log.Println("error marshling game state: ", err)
				return
			}
			conn.WriteMessage(websocket.TextMessage, mapBytes)
			uuid := uuid.New().String()
			playerPositions.Connections[conn] = PlayerID(uuid)
			playerPositions.Positions[PlayerID(uuid)] = maze.Start
			playerPositionsBytes, err := json.Marshal(map[string]interface{}{"playerId": uuid, "position": maze.Start})
			hub.conns[conn] = true
			if err != nil {
				log.Println("error marshling player positions: ", err)
				return
			}
			h.broadcast <- playerPositionsBytes
			log.Println("starting read pump")
			go readPump(h, conn)

		case conn := <-h.unregister:
			if _, ok := h.conns[conn]; ok {
				playerId := playerPositions.Connections[conn]
				delete(playerPositions.Positions, playerId)
				delete(playerPositions.Connections, conn)
				delete(h.conns, conn)
				if err := conn.Close(); err != nil {
					log.Printf("Error Closing connection for player %s: %v", playerId, err)
				}
				deletePlayer, err := json.Marshal(map[string]string{"delete": string(playerId)})
				if err != nil {
					log.Println("error marshling player positions: ", err)
				} else {
					h.broadcast <- deletePlayer
				}
			}
		case message := <-h.broadcast:
			fmt.Println("broadcasting message ", string(message))
			fmt.Println("length of clients, ", len(hub.conns))
			for conn := range h.conns {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					h.unregister <- conn
					log.Println("error writing to websocket: ", err)
				}
			}
		}
	}
}
