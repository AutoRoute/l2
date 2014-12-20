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
			log.Fatal("Failure to read from source", source, err.Error())
		}
		err = destination.WritePacket(p)
		if err != nil {
			log.Fatal("Failure to write to", destination, err.Error())
		}

	}
}

func main() {
	source := "wlp3s0"
	flag.StringVar(&source, "source", source, "source address to listen to")
	mac := "00:24:d7:3e:71:b4"
	flag.StringVar(&mac, "mac", mac, "mac address to use")
	destination := "wlan0"
	flag.StringVar(&destination, "destination", destination, "new device to create")

	macbyte := MAC(mac).ToBytes()

	flag.Parse()

	tap, err := NewTapDevice(mac, destination)
	if err != nil {
		log.Fatal(err)
	}
	defer tap.Close()

	eth, err := ConnectEthDevice(source)
	if err != nil {
		log.Fatal(err)
	}
	filtered_eth := FilterPacket{macbyte, eth}

	go SendPackets(PacketLogger{tap}, eth)
	go SendPackets(PacketLogger{filtered_eth}, tap)

	fmt.Scanln()
}
