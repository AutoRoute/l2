package main

import "testing"

func TestEthToTap(t *testing.T) {
	dev_name := "test0"
	dev_mac := "00:11:22:33:44:55:66"
	dest_mac := "ff:11:22:33:44:55:66"
	tap, err := NewTapDevice(dev_mac, dev_name)
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()

	eth_tap, err := ConnectEthDevice(dev_name)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 100)
	packet := NewEthPacket(MAC(dest_mac).ToBytes(), MAC(dev_mac).ToBytes(), 1, data)
	err = eth_tap.WritePacket(packet)
	if err != nil {
		t.Fatal(err)
	}

	newpacket, err := tap.ReadPacket()
	if err != nil {
		t.Fatal(err)
	}

	if len(newpacket) != len(packet) {
		t.Fatal("Packet length changed ", len(packet), " vs ", len(newpacket))
	}
}

func TestTapToEth(t *testing.T) {
	dev_name := "test1"
	dev_mac := "00:11:22:33:44:55:66"
	dest_mac := "ff:11:22:33:44:55:66"
	tap, err := NewTapDevice(dev_mac, dev_name)
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()

	eth_tap, err := ConnectEthDevice(dev_name)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 100)
	packet := NewEthPacket(MAC(dest_mac).ToBytes(), MAC(dev_mac).ToBytes(), 1, data)
	err = tap.WritePacket(packet)
	if err != nil {
		t.Fatal(err)
	}

	newpacket, err := eth_tap.ReadPacket()
	if err != nil {
		t.Fatal(err)
	}

	if len(newpacket) != len(packet) {
		t.Fatal("Packet length changed ", len(packet), " vs ", len(newpacket))
	}
}
