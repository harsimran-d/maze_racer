package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func readPump(hub *Hub, conn *websocket.Conn) {
	defer func() {
		hub.unregister <- conn
	}()
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
			if _, ok := maze.Finishers[string(playerId)]; ok {
				log.Println("game finisher tried to move")
				continue
			}
			oldPosition := playerPositions.Positions[playerId]
			newPosition, isNew := maze.NewPosition(oldPosition, string(moveBytes))
			if isNew {
				log.Println(newPosition)
				if maze.Finish == newPosition {
					stamp := time.Now()
					maze.Finishers[string(playerId)] = stamp
					finisher, err := json.Marshal(map[string]interface{}{"finisher": playerId, "time": stamp})
					if err != nil {
						log.Println("could not marshal finisher: ", playerId)
					} else {
						hub.broadcast <- finisher
					}
				}
				playerPositions.Positions[playerId] = newPosition
				newMoveBytes, err := json.Marshal(map[string]interface{}{"playerId": playerId, "position": newPosition})
				if err != nil {
					log.Println("could not marshal new position: ", newPosition)
				} else {
					hub.broadcast <- newMoveBytes
				}
			} else {
				log.Println("no new position for you: ", playerId)
			}

		} else {
			log.Printf("unsupported message type: %d", messageType)
		}
	}
}
