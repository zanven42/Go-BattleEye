package BattleEye

import (
	"encoding/binary"
	"errors"
)

func buildHeader(Checksum uint32) []byte {
	Check := make([]byte, 4) // should reduce allocations when i benchmark this shit
	binary.LittleEndian.PutUint32(Check, Checksum)
	// build header and return it.
	return append([]byte{}, 'B', 'E', Check[0], Check[1], Check[2], Check[3])
}

//This Takes a Constructed command or login packet and wraps it with the header and assigned a checksum
func buildPacket(data []byte, PacketType byte) []byte {
	data = append([]byte{0xFF, PacketType}, data...)
	checksum := makeChecksum(data)
	header := buildHeader(checksum)
	return append(header, data...)
}

func buildConnectionPacket(pass string) []byte {
	return buildPacket([]byte(pass), packetType.LOGIN)
}

func buildCommandPacket(command []byte, sequence uint8) []byte {
	return buildPacket(append([]byte{sequence}, command...), packetType.COMMAND)
}

// Note sure if this command packet heartbeat needs to keep track of sequence or not. but thanks battle eye ill presume
// it does since it asks for a 2 byte empty command
func buildHeartBeatPacket(sequence uint8) []byte {
	return buildPacket([]byte{sequence}, packetType.COMMAND)
}

func buildMessageAckPacket(sequence uint8) []byte {
	return buildPacket([]byte{sequence}, packetType.SERVER_MESSAGE)
}

//This function takes in a data packet and returns the type of packet recieved
func responseType(data []byte) (byte, error) {
	if len(data) < 8 {
		return 0, errors.New("Error Packet length too small")
	}
	// 7th element will be the first element after the header which will be the packet type
	return data[7], nil
}