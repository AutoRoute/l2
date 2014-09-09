package main

import (
	"encoding/binary"
	"encoding/hex"
	"log"
)

type EthPacket []byte

func NewEthPacket(destination, source []byte, protocol uint16, data []byte) []byte {
	packet := make([]byte, 12+2+len(data))
	copy(packet[0:6], destination)
	copy(packet[6:12], source)
	binary.BigEndian.PutUint16(packet[12:14], protocol)
	copy(packet[14:], data)
	return packet
}

func (data EthPacket) Destination() []byte {
	return data[0:6]
}

func (data EthPacket) Source() []byte {
	return data[6:12]
}

func (data EthPacket) Type() []byte {
	protocol := int(binary.BigEndian.Uint16(data[12:14]))
	if protocol < 1536 {
		log.Print("Packet protocol is ", hex.EncodeToString(data[12:14]))
	}
	return data[12:14]
}

func (data EthPacket) TypeInt() int {
	protocol := int(binary.BigEndian.Uint16(data[12:14]))
	if protocol < 1536 {
		log.Print("Packet protocol is ", hex.EncodeToString(data[12:14]))
	}
	return protocol
}
