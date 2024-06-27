package packet

import (
	"bytes"
	"encoding/binary"
)

type PingPacket struct {
	Time       int64
	ServerIp   *string
	ServerPort *uint16
}

func DecodePingPacket(packet VoicePacket) (PingPacket, error) {
	buffer := bytes.NewBuffer(packet.Data)

	var time int64
	if err := binary.Read(buffer, binary.BigEndian, &time); err != nil {
		return PingPacket{}, err
	}

	if buffer.Available() == 0 {
		return PingPacket{time, nil, nil}, nil
	}

	var serverIpLen int16
	if err := binary.Read(buffer, binary.BigEndian, &serverIpLen); err != nil {
		return PingPacket{}, err
	}

	serverIpBytes := make([]byte, serverIpLen)
	if _, err := buffer.Read(serverIpBytes); err != nil {
		return PingPacket{}, err
	}

	serverIp := string(serverIpBytes)

	var serverPort uint16
	if err := binary.Read(buffer, binary.BigEndian, &serverPort); err != nil {
		return PingPacket{}, err
	}

	return PingPacket{time, &serverIp, &serverPort}, nil
}
