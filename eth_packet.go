package main

import (
	"encoding/binary"
)

type EthPacket []byte

func NewEthPacket(destination, source []byte, protocol uint16, data []byte) []byte {
	p := make([]byte, 12+2+len(data))
	copy(p[0:6], destination)
	copy(p[6:12], source)
	binary.BigEndian.PutUint16(p[12:14], protocol)
	copy(p[14:], data)
	return p
}

func (p EthPacket) Destination() []byte {
	return p[0:6]
}

func (p EthPacket) Source() []byte {
	return p[6:12]
}

func (p EthPacket) Type() uint16 {
	return binary.BigEndian.Uint16(p[12:14])
}
