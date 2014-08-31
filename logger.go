package main

import (
	"encoding/hex"
	"errors"
	"log"
)

type PacketLogger struct {
	d PacketDevice
}

func (l PacketLogger) ReadPacket() ([]byte, error) {
	for {
		p, err := l.d.ReadPacket()
		if err == nil {
			PrintPacket(l.d.Name(), p)
		}
		return p, err
	}
}

func (l PacketLogger) WritePacket([]byte) error {
	return errors.New("This interface cannot write")
}

func (l PacketLogger) Name() string {
	return "Logger: " + l.d.Name()
}

func PrintPacket(name string, data []byte) {
	E := EthPacket(data)
	log.Printf("%s: %s->%s %d bytes protocol %s",
		name,
		hex.EncodeToString(E.Source()),
		hex.EncodeToString(E.Destination()),
		len(data),
		hex.EncodeToString(E.Type()))
}
