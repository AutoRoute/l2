package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// FilterPacket is a PacketReader which only allows through packets which match the list of
// packets is is supplied with.
type FilterPacket struct {
	mac    [][]byte
	device PacketReader
}

func NewFilterPacket(dev PacketReader, mac ...[]byte) *FilterPacket {
	return &FilterPacket{mac, dev}
}

func (f FilterPacket) ReadPacket() ([]byte, error) {
	for {
		p, err := f.device.ReadPacket()
		if err != nil {
			return p, err
		}
		for _, mac := range f.mac {
			if bytes.Equal(EthPacket(p).Destination(), mac) {
				return p, nil
			}
		}
	}
}

func (f FilterPacket) String() string {
	s := "FilterPacket{" + fmt.Sprint(f.device)
	for _, mac := range f.mac {
		s += ", " + hex.EncodeToString(mac)
	}
	s += "}"
	return s
}
