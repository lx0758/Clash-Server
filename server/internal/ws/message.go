package ws

import "encoding/json"

type MessageType string

const (
	TypeTraffic     MessageType = "traffic"
	TypeConnections MessageType = "connections"
	TypeLogs        MessageType = "logs"
	TypeCoreStatus  MessageType = "core_status"
	TypeMemory      MessageType = "memory"
)

type Message struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type TrafficData struct {
	Up   int64 `json:"up"`
	Down int64 `json:"down"`
}

type CoreStatusData struct {
	Running bool   `json:"running"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ConnectionsData struct {
	Connections   []interface{} `json:"connections"`
	DownloadTotal int64         `json:"downloadTotal"`
	UploadTotal   int64         `json:"uploadTotal"`
}

type LogData struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type MemoryData struct {
	Inuse   int64 `json:"inuse"`
	Oslimit int64 `json:"oslimit"`
}

type SubscribeRequest struct {
	Action string   `json:"action"`
	Types  []string `json:"types"`
}
