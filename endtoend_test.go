package main

import (
    "testing"
    "bytes"
    "time"
    "log"
)

func waitForPacket(target []byte, dev PacketDevice) bool {
    done := make(chan struct{})
    found := make(chan struct{})
    go func() {
        for {
            select{
            case <- done:
                return
            default:
            }

            pack, err := dev.ReadPacket()
            if err != nil {
                log.Fatal(err)
            }
            if bytes.Equal(target, pack) {
                found <- struct{}{}
                return
            }
            PrintPacket("desired", target)
            PrintPacket("found", pack)
        }
    }()
    
    timeout := time.After(time.Second)

    select {
        case <- timeout:
            done <- struct{}{}
            return false
        case <- found:
            return true
    }
}

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

    waitForPacket(packet, tap)
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

    waitForPacket(packet, eth_tap)
}
