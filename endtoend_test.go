package l2

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

var (
	dev_name = "test"
	dev_mac  = "00:11:22:33:44:55"
	dest_mac = "ff:11:22:33:44:55"
)

func waitForFrame(target []byte, dev FrameReader) bool {
	timeout := time.After(time.Second)
	for {
		// Check if we're out of time
		select {
		case <-timeout:
			return false
		default:
		}

		frame, err := dev.ReadFrame()
		if err != nil {
			log.Print(err)
			return false
		}
		if bytes.Equal(target, frame) {
			return true
		}
		PrintFrame("desired", target)
		PrintFrame("found", frame)
	}
}

func NewDevices() (*TapDevice, FrameReadWriter, error) {
	tap, err := NewTapDevice(dev_mac, dev_name)
	if err != nil {
		return nil, nil, err
	}

	eth, err := ConnectExistingDevice(dev_name)
	if err != nil {
		tap.Close()
		return nil, nil, err
	}
	return tap, eth, nil
}

func NewFrame(dest, src string) EthFrame {
	data := make([]byte, 100)
	return NewEthFrame(macToBytesOrDie(dest), macToBytesOrDie(src), 1, data)
}

func TestEthToTap(t *testing.T) {
	tap, eth, err := NewDevices()
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()
	frame := NewFrame(dest_mac, dev_mac)
	err = eth.WriteFrame(frame)
	if err != nil {
		t.Fatal(err)
	}
	waitForFrame(frame, tap)
}

func TestTapToEth(t *testing.T) {
	tap, eth, err := NewDevices()
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()
	frame := NewFrame(dest_mac, dev_mac)
	err = tap.WriteFrame(frame)
	if err != nil {
		t.Fatal(err)
	}
	waitForFrame(frame, eth)
}
func TestPrinting(t *testing.T) {
	tap, eth, err := NewDevices()
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()
	if !strings.Contains(fmt.Sprint(tap), "TapDevice") {
		t.Fatal("Missing TapDevice from", tap)
	}
	if !strings.Contains(fmt.Sprint(eth), "existingDevice") {
		t.Fatal("Missing existingDevice from", eth)
	}
}
