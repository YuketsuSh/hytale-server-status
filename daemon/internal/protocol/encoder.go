package protocol

import (
	"encoding/binary"
	"fmt"
)

func WriteVarInt(buf []byte, value uint32) int {
	i := 0
	for value >= 0x80 {
		buf[i] = byte(value&0x7F | 0x80)
		value >>= 7
		i++
	}
	buf[i] = byte(value)
	return i + 1
}

func ReadVarInt(buf []byte) (uint32, int) {
	var result uint32
	var shift uint
	i := 0

	for {
		b := buf[i]
		i++
		result |= uint32(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift > 35 {
			panic("varint too large")
		}
	}

	return result, i
}

func WriteString(buf []byte, s string) int {
	length := len(s)
	offset := WriteVarInt(buf, uint32(length))
	copy(buf[offset:], s)
	return offset + length
}

func ReadString(buf []byte) (string, int) {
	length, offset := ReadVarInt(buf)
	end := offset + int(length)
	if end > len(buf) {
		return "", offset
	}
	return string(buf[offset:end]), end
}

func SerializeConnectPacket(packet *ConnectPacket) ([]byte, error) {
	totalSize := 0

	totalSize += 64 + 1 + 16

	totalSize += len(packet.Language) + 5
	totalSize += len(packet.IdentityToken) + 5
	totalSize += len(packet.Username) + 5
	totalSize += len(packet.ReferralData) + 5
	totalSize += len(packet.ReferralSource) + 5

	buf := make([]byte, totalSize+8)
	offset := 0

	offset += 4

	binary.LittleEndian.PutUint32(buf[offset:], PACKET_CONNECT)
	offset += 4

	if len(packet.ProtocolHash) != 64 {
		return nil, fmt.Errorf("protocol hash must be 64 bytes, got %d", len(packet.ProtocolHash))
	}
	copy(buf[offset:], []byte(packet.ProtocolHash))
	offset += 64

	buf[offset] = byte(packet.ClientType)
	offset += 1

	offset += WriteString(buf[offset:], packet.Language)

	offset += WriteString(buf[offset:], packet.IdentityToken)

	copy(buf[offset:], packet.UUID[:])
	offset += 16

	offset += WriteString(buf[offset:], packet.Username)

	if len(packet.ReferralData) > 0 {
		offset += WriteVarInt(buf[offset:], uint32(len(packet.ReferralData)))
		copy(buf[offset:], packet.ReferralData)
		offset += len(packet.ReferralData)
	} else {
		offset += WriteVarInt(buf[offset:], 0)
	}

	offset += WriteString(buf[offset:], packet.ReferralSource)

	packetLength := offset - 4
	binary.LittleEndian.PutUint32(buf[0:], uint32(packetLength))

	return buf[:offset], nil
}
