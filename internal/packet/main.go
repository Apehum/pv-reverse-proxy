package packet

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
)

type VoicePacket struct {
	Type   byte
	Secret string
	Data   []byte
}

func DecodePacket(packet []byte) (VoicePacket, error) {
	buffer := bytes.NewBuffer(packet)

	var magicNumber int32
	if err := binary.Read(buffer, binary.LittleEndian, &magicNumber); err != nil {
		return VoicePacket{}, err
	}

	var packetType byte
	if p, err := buffer.ReadByte(); err != nil {
		return VoicePacket{}, err
	} else {
		packetType = p
	}

	secretBytes := make([]byte, 16)
	if _, err := buffer.Read(secretBytes); err != nil {
		return VoicePacket{}, err
	}

	secret := uuid.UUID(secretBytes)

	// skip packet time, it's useless and not used
	_, _ = buffer.Read(make([]byte, 8))

	data := buffer.Bytes()

	return VoicePacket{packetType, secret.String(), data}, nil
}
