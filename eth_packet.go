package main

import (
	"encoding/binary"
	"encoding/hex"
	"log"
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

func (p EthPacket) checkType() {
	protocol := int(binary.BigEndian.Uint16(p[12:14]))
	if protocol < 1536 {
		log.Print("This packet appears to not be a 802.3ad packet which is unsupported. The protocol is ", hex.EncodeToString(p[12:14]))
	}
}
func (p EthPacket) Type() []byte {
	p.checkType()
	return p[12:14]
}

func (p EthPacket) TypeInt() int {
	p.checkType()
	return int(binary.BigEndian.Uint16(p[12:14]))
}
