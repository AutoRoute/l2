package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type FilterPacket struct {
	mac    []byte
	device PacketReader
}

func (f FilterPacket) ReadPacket() ([]byte, error) {
	broadcast := MAC("ff:ff:ff:ff:ff:ff").ToBytes()
	for {
		p, err := f.device.ReadPacket()
		if err != nil {
			return p, err
		}
		if bytes.Equal(EthPacket(p).Destination(), f.mac) {
			return p, nil
		}
		if bytes.Equal(EthPacket(p).Destination(), broadcast) {
			return p, nil
		}
	}
}

func (f FilterPacket) String() string {
	return "FilterPacket{" + hex.EncodeToString(f.mac) + ", " + fmt.Sprint(f.device) + "}"
}
