package main

import (
	"bytes"
)

type FilterPacket struct {
	mac    []byte
	reader PacketDevice
}

func (f FilterPacket) ReadPacket() []byte {
	for {
		p := f.reader.ReadPacket()
		if bytes.Equal(EthPacket(p).Destination(), f.mac) {
			return p
		}
		broadcast := MAC("ff:ff:ff:ff:ff:ff").ToBytes()
		if bytes.Equal(EthPacket(p).Destination(), broadcast) {
			return p
		}
	}
}

func (f FilterPacket) WritePacket(data []byte) {}

func (f FilterPacket) Name() string {
	return "wrapped: " + f.reader.Name()
}
