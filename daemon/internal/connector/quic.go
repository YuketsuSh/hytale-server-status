package connector

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	_ "net"
	"time"

	"daemon/internal/protocol"

	"github.com/quic-go/quic-go"
)

type HytaleConnector struct {
	timeout time.Duration
}

func NewHytaleConnector(timeout time.Duration) *HytaleConnector {
	return &HytaleConnector{
		timeout: timeout,
	}
}

func (c *HytaleConnector) Connect(ctx context.Context, host string, port int) (quic.Connection, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"hytale/1"},
	}

	quicConfig := &quic.Config{
		HandshakeIdleTimeout: c.timeout,
		MaxIdleTimeout:       c.timeout,
		KeepAlivePeriod:      c.timeout / 2,
	}

	dialCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	conn, err := quic.DialAddr(dialCtx, address, tlsConfig, quicConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to establish QUIC connection to %s: %w", address, err)
	}

	return conn, nil
}

func (c *HytaleConnector) SendPacket(conn quic.Connection, packetID uint32, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return fmt.Errorf("failed to open QUIC stream: %w", err)
	}
	defer stream.Close()

	packetLength := uint32(len(data) + 4)

	lengthBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBuf, packetLength)
	if _, err := stream.Write(lengthBuf); err != nil {
		return fmt.Errorf("failed to write packet length: %w", err)
	}

	idBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(idBuf, packetID)
	if _, err := stream.Write(idBuf); err != nil {
		return fmt.Errorf("failed to write packet ID: %w", err)
	}

	if _, err := stream.Write(data); err != nil {
		return fmt.Errorf("failed to write packet data: %w", err)
	}

	return nil
}

func (c *HytaleConnector) ReceivePacket(conn quic.Connection) (*protocol.Packet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	stream, err := conn.AcceptStream(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to accept QUIC stream: %w", err)
	}
	defer stream.Close()

	lengthBuf := make([]byte, 4)
	if _, err := stream.Read(lengthBuf); err != nil {
		return nil, fmt.Errorf("failed to read packet length: %w", err)
	}

	length := binary.LittleEndian.Uint32(lengthBuf)
	if length > protocol.MAX_PACKET_SIZE {
		return nil, fmt.Errorf("packet too large: %d bytes", length)
	}
	if length < 4 {
		return nil, fmt.Errorf("invalid packet length: %d", length)
	}

	idBuf := make([]byte, 4)
	if _, err := stream.Read(idBuf); err != nil {
		return nil, fmt.Errorf("failed to read packet ID: %w", err)
	}

	packetID := binary.LittleEndian.Uint32(idBuf)

	payloadSize := length - 4
	payload := make([]byte, payloadSize)
	if _, err := stream.Read(payload); err != nil {
		return nil, fmt.Errorf("failed to read packet payload: %w", err)
	}

	return &protocol.Packet{
		ID:      packetID,
		Length:  length,
		Payload: payload,
	}, nil
}

func (c *HytaleConnector) QueryServerStatus(host string, port int) (*protocol.ServerStatus, time.Duration, error) {
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	conn, err := c.Connect(ctx, host, port)
	if err != nil {
		return &protocol.ServerStatus{
			Address:      fmt.Sprintf("%s:%d", host, port),
			Online:       false,
			ErrorMessage: err.Error(),
			Timestamp:    time.Now().Unix(),
		}, 0, err
	}
	defer conn.CloseWithError(0, "status query complete")

	connectPacket := protocol.NewConnectPacket("HytaleStatusDaemon")
	packetData, err := protocol.SerializeConnectPacket(connectPacket)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to serialize connect packet: %w", err)
	}

	if err := c.SendPacket(conn, protocol.PACKET_CONNECT, packetData[4:]); err != nil {
		return nil, 0, fmt.Errorf("failed to send connect packet: %w", err)
	}

	packet, err := c.ReceivePacket(conn)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to receive response: %w", err)
	}

	latency := time.Since(startTime)

	if packet.ID == protocol.PACKET_STATUS {
		// TODO: Parse Status packet payload
		status := &protocol.StatusPacket{
			PlayerCount:   0,
			MaxPlayers:    0,
			ServerName:    "",
			MOTD:          "",
			ServerVersion: "",
		}

		serverStatus := status.ToServerStatus(fmt.Sprintf("%s:%d", host, port), latency)
		return serverStatus, latency, nil
	}

	return &protocol.ServerStatus{
		Address:      fmt.Sprintf("%s:%d", host, port),
		Online:       false,
		ErrorMessage: fmt.Sprintf("unexpected packet ID: %d", packet.ID),
		Timestamp:    time.Now().Unix(),
	}, latency, fmt.Errorf("unexpected packet ID: %d", packet.ID)
}
