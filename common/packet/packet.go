package packet

import (
    "bufio"
    "bytes"
    "encoding/binary"
    "errors"
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
	CMD_REGISTER
	CMD_LOGIN
	CMD_SEARCH
	CMD_USER
	CMD_MSG

	// Server Responses
	SVR_NOTICE
	SVR_MSGTO
	SVR_VERSION
	SVR_LOGIN
	SVR_REGISTER
	SVR_ERROR
	SVR_USER
	SVR_MSG
	SVR_BLOCK
)

const HEADER_BYTE byte = '\xde'
const MAX_NAME_LENGTH int = 65535

const maxPacketLength = 1<<20 // 1 MiB safety cap

func PacketRead (br *bufio.Reader) ([]byte, error) {
    _, err := br.ReadBytes(HEADER_BYTE)
	if err != nil {
		return nil, err
	}
	// Get packet length field
	packet := make([]byte, 4)
	read_bytes := 0
	for read_bytes < 4 {
		tmp := make([]byte, 4)
		nread, err := br.Read(tmp)
		if err != nil {
			return nil, err
		}
		copy(packet[read_bytes:], tmp[:nread])
		read_bytes += nread
	}
	// Get rest of the packet
    packetLen := int(binary.BigEndian.Uint32(packet))
    if packetLen <= 0 || packetLen > maxPacketLength {
        return nil, errors.New("invalid packet length")
    }
	packet = make([]byte, packetLen)
	read_bytes = 0
	for read_bytes < packetLen {
		tmp := make([]byte, packetLen)
		nread, err := br.Read(tmp)
		if err != nil {
			return nil, err
		}
		copy(packet[read_bytes:], tmp[:nread])
		read_bytes += nread
	}
	return packet, nil
}

func PacketForm(packetType byte, payload []byte) []byte {
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
