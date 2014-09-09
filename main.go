package main

import (
	"flag"
	"fmt"
	"log"
)

func SendPackets(source PacketReaderDevice, destination PacketWriterDevice) {
	for {
		p, err := source.ReadPacket()
		if err != nil {
			log.Fatal("Failure to read from source", source.Name(), err.Error())
		}
		err = destination.WritePacket(p)
		if err != nil {
			log.Fatal("Failure to write to", destination.Name(), err.Error())
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

	fd, err := NewTapDevice(mac, destination)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	raw_eth, err := ConnectEthDevice(source)
	if err != nil {
		log.Fatal(err)
	}
	eth := FilterPacket{macbyte, raw_eth}

	go SendPackets(PacketLogger{fd}, eth)
	go SendPackets(PacketLogger{eth}, fd)

	fmt.Scanln()
}
