package main

import (
	"encoding/hex"
	"log"
)

type PacketLogger struct {
	d PacketDevice
}

func (l PacketLogger) ReadPacket() []byte {
	for {
		p := l.d.ReadPacket()
		PrintPacket(l.d.Name(), p)
		return p
	}
}

func (l PacketLogger) WritePacket([]byte) {
	panic("ERROR")
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
