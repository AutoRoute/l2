package main

import (
	"bytes"
	"log"
	"testing"
	"time"
)

var (
	dev_name = "test"
	dev_mac  = "00:11:22:33:44:55"
	dest_mac = "ff:11:22:33:44:55"
)

func waitForPacket(target []byte, dev PacketReader) bool {
	timeout := time.After(time.Second)
	for {
		// Check if we're out of time
		select {
		case <-timeout:
			return false
		default:
		}

		pack, err := dev.ReadPacket()
		if err != nil {
			log.Print(err)
			return false
		}
		if bytes.Equal(target, pack) {
			return true
		}
		PrintPacket("desired", target)
		PrintPacket("found", pack)
	}
}

func NewDevices() (*TapDevice, *EthDevice, error) {
	tap, err := NewTapDevice(dev_mac, dev_name)
	if err != nil {
		return nil, nil, err
	}

	eth, err := ConnectEthDevice(dev_name)
	if err != nil {
		tap.Close()
		return nil, nil, err
	}
	return tap, eth, nil
}

func NewPacket(dest, src string) []byte {
	data := make([]byte, 100)
	return NewEthPacket(MacToBytesOrDie(dest), MacToBytesOrDie(src), 1, data)
}

func TestEthToTap(t *testing.T) {
	tap, eth, err := NewDevices()
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()
	packet := NewPacket(dest_mac, dev_mac)
	err = eth.WritePacket(packet)
	if err != nil {
		t.Fatal(err)
	}
	waitForPacket(packet, tap)
}

func TestTapToEth(t *testing.T) {
	tap, eth, err := NewDevices()
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()
	packet := NewPacket(dest_mac, dev_mac)
	err = tap.WritePacket(packet)
	if err != nil {
		t.Fatal(err)
	}
	waitForPacket(packet, eth)
}
