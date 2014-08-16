package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
)

func SendPackets(source PacketReader, destination PacketWriter) {
	for {
		p := source.ReadPacket()
		destination.WritePacket(p)
	}
}

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
		broadcast, err := hex.DecodeString("ffffffffffff")
		if err != nil {
			log.Fatal(err)
		}
		if bytes.Equal(EthPacket(p).Destination(), broadcast) {
			return p
		}
	}
}

func (f FilterPacket) WritePacket(data []byte) {}

func (f FilterPacket) Name() string {
	return "wrapped: " + f.reader.Name()
}

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

func main() {
	mac := "00:24:d7:3e:71:b4"
	macbyte := MAC(mac).ToBytes()

	fd := NewTapDevice(mac)
	defer fd.Close()

	eth := FilterPacket{macbyte, ConnectEthDevice("wlp3s0")}

	go SendPackets(PacketLogger{fd}, eth)
	go SendPackets(PacketLogger{eth}, fd)

	fmt.Scanln()
}
