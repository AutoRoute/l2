// Program l2forward is a simple binary to foward networking devies at the layer two level.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"
)

// These two interfaces represent the core abstractions in l2forward, with
// almost everything working on PacketReaders. This is not implemented as a
// io.Reader because you cannot arbitrarily slice up layer 2 ethernet packets
// and expect things to keep working
type PacketReader interface {
	ReadPacket() ([]byte, error)
}

type PacketWriter interface {
	WritePacket([]byte) error
}

// Basic utility function to take a string and turn it into a mac address. Will
// die if the string is not valid.
func MacToBytesOrDie(m string) []byte {
	b, err := hex.DecodeString(strings.Replace(m, ":", "", -1))
	if err != nil {
		log.Fatal(err)
	}
	if len(b) != 6 {
		log.Fatal("Expected mac of length 6 bytes got ", len(b))
	}
	return b
}

// Local equivalent of io.Copy, will shove packets from a PacketReader
// into a PacketWriter
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
