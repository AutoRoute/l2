package main

import (
	"encoding/hex"
	"fmt"
	"log"
)

type PacketLogger struct {
	d PacketReader
}

func (l PacketLogger) ReadPacket() ([]byte, error) {
	for {
		p, err := l.d.ReadPacket()
		if err == nil {
			PrintPacket(fmt.Sprint(l.d), p)
		} else {
			log.Print("Err reading packet:", err)
		}
		return p, err
	}
}

func (l PacketLogger) String() string {
	return "Logger{" + fmt.Sprint(l.d) + "}"
}

func PrintPacket(name string, data []byte) {
	E := EthPacket(data)
	log.Printf("%s: %s->%s %d bytes protocol %d",
		name,
		hex.EncodeToString(E.Source()),
		hex.EncodeToString(E.Destination()),
		len(data),
		E.Type())
}
