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
	dev := flag.String("dev", "wlan0", "Device to create/attach to")
	mac := flag.String("mac", "e8:b1:fc:07:fa:3f", "mac address to use")
	broadcast := flag.String("broadcast", "", "Address to listen on (mutually exclusive with -connect)")
	connect := flag.String("connect", "", "Address to connect to (mutually exclusive with -broadcast)")
	flag.Parse()

	if len(*broadcast) == 0 && len(*connect) == 0 {
		log.Fatal("Must specify broadcast or connect")
	}

	if len(*broadcast) != 0 && len(*connect) != 0 {
		log.Fatal("Cannot specify broadcast and connect")
	}

	macbyte := MacToBytesOrDie(*mac)
	macbroad := MacToBytesOrDie("ff:ff:ff:ff:ff:ff")

	if len(*broadcast) != 0 {
		eth, err := ConnectEthDevice(*dev)
		if err != nil {
			log.Fatal(err)
		}
		filtered_eth := NewFilterPacket(eth, macbroad, macbyte)
		ln, err := NewListener(*broadcast)
		if err != nil {
			log.Fatal(err)
		}
		go SendPackets(PacketLogger{ln}, eth)
		go SendPackets(PacketLogger{filtered_eth}, ln)
	} else {
		tap, err := NewTapDevice(*mac, *dev)
		if err != nil {
			log.Fatal(err)
		}
		defer tap.Close()
		c, err := NewDialer(*connect)
		if err != nil {
			log.Fatal(err)
		}
		go SendPackets(PacketLogger{tap}, c)
		go SendPackets(PacketLogger{c}, tap)
	}
	fmt.Scanln()
}
