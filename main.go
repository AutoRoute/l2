package main

import (
	"flag"
	"fmt"
	"log"
)

func SendPackets(source PacketReader, destination PacketWriter) {
	for {
		p, err := source.ReadPacket()
		if err != nil {
			log.Fatal("Failure to read from source", source, err)
		}
		if destination.WritePacket(p) != nil {
			log.Fatal("Failure to write to", destination, err)
		}
	}
}

func main() {
    source := flag.String("source", "wlp3s0", "source address to listen to")
    mac := flag.String("mac", "00:24:d7:3e:71:b4", "mac address to use")
    destination := flag.String("destination", "wlan0", "new device to create")
	flag.Parse()

	macbyte := MacToBytesOrDie(*mac)
	broadcast := MacToBytesOrDie("ff:ff:ff:ff:ff:ff")

	tap, err := NewTapDevice(*mac, *destination)
	if err != nil {
		log.Fatal(err)
	}
	defer tap.Close()

	eth, err := ConnectEthDevice(*source)
	if err != nil {
		log.Fatal(err)
	}
	filtered_eth := NewFilterPacket(eth, broadcast, macbyte)

	go SendPackets(PacketLogger{tap}, eth)
	go SendPackets(PacketLogger{filtered_eth}, tap)

	fmt.Scanln()
}
