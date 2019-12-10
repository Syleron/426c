package main

import (
	"bytes"
	"encoding/binary"
)

// Packet Structure [header (1)] [packet length (4)] [type (1)] [payload]
// Packet Length = Total bytes - header size
const (
	// Client Commands
	CMD_MSGALL = iota
	CMD_MSGTO
	CMD_IDENT
	CMD_WHO
	CMD_JOINCLUB
	CMD_LEAVECLUB

	// Server Responses
	SVR_NOTICE
	SVR_MSG
	SVR_VERSION
)

const HEADER_BYTE byte = '\xde'
const MAX_NAME_LENGTH int = 65535

func packetForm(packetType byte, payload []byte) []byte {
	var buf bytes.Buffer
	buf.Write([]byte{
		HEADER_BYTE,
	})
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(payload)+1)); err != nil {
		return nil
	}
	buf.Write([]byte{
		packetType,
	})
	buf.Write(payload)
	return buf.Bytes()
}
