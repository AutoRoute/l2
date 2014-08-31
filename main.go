package main

import (
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
	mac := "00:24:d7:3e:71:b4"
	macbyte := MAC(mac).ToBytes()

	fd, err := NewTapDevice(mac, "wlan0")
    if err != nil {
        log.Fatal(err)
    }
	defer fd.Close()

	raw_eth, err := ConnectEthDevice("wlp3s0")
    if err != nil {
        log.Fatal(err)
    }
	eth := FilterPacket{macbyte, raw_eth}

	go SendPackets(PacketLogger{fd}, eth)
	go SendPackets(PacketLogger{eth}, fd)

	fmt.Scanln()
}
