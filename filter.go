package main

import (
	"bytes"
)

type FilterPacket struct {
	mac    []byte
	device PacketDevice
}

func (f FilterPacket) ReadPacket() ([]byte, error) {
	for {
		p, err := f.device.ReadPacket()
		if err != nil {
			return nil, err
		}
		if bytes.Equal(EthPacket(p).Destination(), f.mac) {
			return p, nil
		}
		broadcast := MAC("ff:ff:ff:ff:ff:ff").ToBytes()
		if bytes.Equal(EthPacket(p).Destination(), broadcast) {
			return p, nil
		}
	}
}

func (f FilterPacket) WritePacket(data []byte) error {
	return f.device.WritePacket(data)
}

func (f FilterPacket) Name() string {
	return "wrapped: " + f.device.Name()
}
