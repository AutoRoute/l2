package main

import (
	"fmt"
)

func SendPackets(source PacketReader, destination PacketWriter) {
	for {
		p := source.ReadPacket()
		destination.WritePacket(p)
	}
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
