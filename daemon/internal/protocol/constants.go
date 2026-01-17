package protocol

const (
	PROTOCOL_HASH    = "6708f121966c1c443f4b0eb525b2f81d0a8dc61f5003a692a8fa157e5e02cea9"
	DEFAULT_PORT     = 5520
	MAX_PACKET_SIZE  = 1677721600
	PROTOCOL_VERSION = 1
	PACKET_COUNT     = 268
	STRUCT_COUNT     = 315
	ENUM_COUNT       = 136
)

const (
	PACKET_CONNECT        = 0
	PACKET_DISCONNECT     = 1
	PACKET_PING           = 2
	PACKET_PONG           = 3
	PACKET_STATUS         = 10
	PACKET_AUTH_TOKEN     = 12
	PACKET_CONNECT_ACCEPT = 14
)

type ClientType byte

const (
	CLIENT_TYPE_GAME   ClientType = 0
	CLIENT_TYPE_EDITOR ClientType = 1
)

type DisconnectType byte

const (
	DISCONNECT_NORMAL DisconnectType = 0
	DISCONNECT_CRASH  DisconnectType = 1
)

type PongType byte

const (
	PONG_TYPE_RAW    PongType = 0
	PONG_TYPE_DIRECT PongType = 1
	PONG_TYPE_TICK   PongType = 2
)

var (
	ErrInvalidPacket    = "invalid packet structure"
	ErrProtocolMismatch = "protocol version mismatch"
	ErrServerTimeout    = "server timeout"
	ErrConnectionFailed = "connection failed"
	ErrInvalidResponse  = "invalid server response"
	ErrPacketTooLarge   = "packet exceeds maximum size"
	ErrInvalidPacketID  = "invalid packet ID"
)
