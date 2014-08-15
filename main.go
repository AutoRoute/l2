package main

import (
	"encoding/hex"
	"fmt"
	"log"
)

func PrintPacketReader(t PacketDevice) {
	for {
		p := t.ReadPacket()
		PrintPacket(t.Name(), p)
	}
}

func PrintPacket(name string, data []byte) {
	log.Printf("%s: %s->%s %d bytes", name,
		hex.EncodeToString(data[0:6]),
		hex.EncodeToString(data[6:12]),
		len(data))
}

func main() {
	fd := NewTapDevice()
	defer fd.Close()
	go PrintPacketReader(fd)
	go PrintPacketReader(ConnectEthDevice("wlp3s0"))
	fmt.Scanln()
}
