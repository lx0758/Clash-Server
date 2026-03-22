package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"clash-server/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hub *ws.Hub

func InitWebSocketHub() *ws.Hub {
	hub = ws.NewHub()
	go hub.Run()
	return hub
}

func WebSocketHandler(c *gin.Context) {
	log.Printf("[WebSocket] Upgrading connection")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] Upgrade error: %v", err)
		return
	}
	log.Printf("[WebSocket] Connection established")

	client := ws.NewClient(hub)
	hub.Register(client)
	defer func() {
		hub.Unregister(client)
		conn.Close()
		log.Printf("[WebSocket] Connection closed")
	}()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var req ws.SubscribeRequest
			if err := json.Unmarshal(message, &req); err == nil {
				var types []ws.MessageType
				for _, t := range req.Types {
					types = append(types, ws.MessageType(t))
				}
				if req.Action == "subscribe" {
					client.Subscribe(types)
					log.Printf("[WebSocket] Client subscribed to: %v", types)
				} else if req.Action == "unsubscribe" {
					client.Unsubscribe(types)
				}
			}
		}
	}()

	for message := range client.GetSend() {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
	conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func BroadcastTraffic(up, down int64) {
	if hub != nil {
		hub.BroadcastToType(ws.TypeTraffic, &ws.TrafficData{Up: up, Down: down})
	}
}

func BroadcastCoreStatus(running bool, version string, errMsg string) {
	if hub != nil {
		hub.BroadcastToType(ws.TypeCoreStatus, &ws.CoreStatusData{
			Running: running,
			Version: version,
			Error:   errMsg,
		})
	}
}
