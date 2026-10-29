package ws

import "encoding/json"

type MessageType string

const (
	TypeTraffic     MessageType = "traffic"
	TypeConnections MessageType = "connections"
	TypeLogs        MessageType = "logs"
	TypeCoreStatus  MessageType = "core_status"
	TypeMemory      MessageType = "memory"
	TypeDelay       MessageType = "delay"
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

// DelayData 承载测延迟过程中单个节点的延迟结果。
// Done 为 true 时表示该组的全部节点测试完成, 此时不携带 Name/Delay。
type DelayData struct {
	Group string `json:"group"`
	Name  string `json:"name,omitempty"`
	Delay int    `json:"delay,omitempty"`
	Done  bool   `json:"done,omitempty"`
	Error string `json:"error,omitempty"`
}

type SubscribeRequest struct {
	Action string   `json:"action"`
	Types  []string `json:"types"`
}
