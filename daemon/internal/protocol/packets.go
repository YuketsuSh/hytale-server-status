package protocol

import (
	_ "encoding/binary"
	"time"
)

type Packet struct {
	ID      uint32
	Length  uint32
	Payload []byte
}

type ConnectPacket struct {
	ProtocolHash   string
	ClientType     ClientType
	Language       string
	IdentityToken  string
	UUID           [16]byte
	Username       string
	ReferralData   []byte
	ReferralSource string
}

type StatusPacket struct {
	PlayerCount   int32
	MaxPlayers    int32
	ServerName    string
	MOTD          string
	ProtocolHash  string
	ServerVersion string
	ExtraData     map[string]interface{}
}

type PingPacket struct {
	ID                  int32
	Time                InstantData
	LastPingValueRaw    int32
	LastPingValueDirect int32
	LastPingValueTick   int32
}

type PongPacket struct {
	ID              int32
	Time            InstantData
	Type            PongType
	PacketQueueSize int16
}

type InstantData struct {
	Seconds     int64
	Nanoseconds int32
}

type ServerStatus struct {
	Address       string `json:"address"`
	Online        bool   `json:"online"`
	MOTD          string `json:"motd,omitempty"`
	ServerVersion string `json:"server_version,omitempty"`
	Players       struct {
		Online int `json:"online"`
		Max    int `json:"max"`
	} `json:"players"`
	LatencyMS    int64  `json:"latency_ms"`
	PacketType   string `json:"packet_type"`
	ErrorMessage string `json:"error_message,omitempty"`
	Timestamp    int64  `json:"timestamp"`
}

func NewConnectPacket(username string) *ConnectPacket {
	return &ConnectPacket{
		ProtocolHash: PROTOCOL_HASH,
		ClientType:   CLIENT_TYPE_GAME,
		Language:     "en-US",
		Username:     username,
		UUID:         [16]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
	}
}

func (sp *StatusPacket) ToServerStatus(address string, latency time.Duration) *ServerStatus {
	status := &ServerStatus{
		Address:       address,
		Online:        true,
		MOTD:          sp.MOTD,
		ServerVersion: sp.ServerVersion,
		LatencyMS:     latency.Milliseconds(),
		PacketType:    "status",
		Timestamp:     time.Now().Unix(),
	}
	status.Players.Online = int(sp.PlayerCount)
	status.Players.Max = int(sp.MaxPlayers)
	return status
}
