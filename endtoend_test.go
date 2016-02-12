package l2

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"testing"
	"time"
)

var (
	dev_name = "test"
	dev_mac  = "00:11:22:33:44:55"
	dest_mac = "ff:11:22:33:44:55"
)

// Waits for a frame that matches a desired one.
// Args:
//	target: The frame we are looking for.
//	dev: The dvice we are waiting on.
// Returns:
//	The number of bytes it read total, error.
func waitForFrame(target []byte, dev FrameReader) (int, error) {
	timeout := time.After(time.Second)
	bytes_read := 0
	for {
		// Check if we're out of time
		select {
		case <-timeout:
			return 0, errors.New("waitForFrame timed out.")
		default:
		}

		frame, err := dev.ReadFrame()
		if err != nil {
			log.Print(err)
			return 0, err
		}
		bytes_read += len(frame)

		if bytes.Equal(target, frame) {
			return bytes_read, nil
		}
		PrintFrame("desired", target)
		PrintFrame("found", frame)
	}
}

// Creates a new connected tap/eth pair.
// Returns:
//	The tap device and the ethernet device, error.
func NewDevices() (FrameReadWriteCloser, FrameReadWriter, error) {
	tap, err := NewTapDevice(dev_mac, dev_name)
	if err != nil {
		return nil, nil, err
	}

	return makeEth(tap)
}

// Creates a new connected tap/eth pair with bandwidth restrictions.
// Args:
//	send_bandwidth: Max bandwidth for outgoing data. (b/s)
//	receive_bandwidth: Max bandwidth for incoming data. (b/s)
// Returns:
//	The tap device and the ethernet device, error.
func NewDevicesWithBandwidth(send_bandwidth,
	receive_bandwidth int) (FrameReadWriteCloser,
	FrameReadWriter, error) {
	tap, err := NewTapDeviceWithLatency(dev_mac, dev_name, send_bandwidth,
		receive_bandwidth)
	if err != nil {
		return nil, nil, err
	}

	return makeEth(tap)
}

// Makes an eth device and connects an existing tap device to it.
// Args:
// 	tap: The tap device.
// Returns:
// 	The tap device and the ethernet device, error.
func makeEth(tap FrameReadWriteCloser) (FrameReadWriteCloser,
	FrameReadWriter, error) {
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

	_, err = waitForFrame(frame, tap)
	if err != nil {
		t.Fatal(err)
	}
}

// Test sending packets with restricted bandwidth.
func TestSendingBandwidth(t *testing.T) {
	tap, _, err := NewDevicesWithBandwidth(10000, 10000)
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()

	frame := NewFrame(dest_mac, dev_mac)

	start_time := time.Now()
	bytes_written := 0
	// Write the frame a lot of times so we get a more accurate idea of bandwidth.
	for i := 0; i < 100; i++ {
		err = tap.WriteFrame(frame)
		if err != nil {
			t.Fatal(err)
		}
		bytes_written += len(frame)
	}
	end_time := time.Now()

	// Check outgoing bandwidth.
	elapsed := float64(end_time.Sub(start_time)) / float64(time.Second)
	bandwidth := float64(bytes_written) / elapsed

	if math.Abs(10000-bandwidth) > 50 {
		t.Fatalf("Got %f b/s of bandwidth, expected 10000.", bandwidth)
	}
}

// Test receiving packets with restricted bandwidth.
func TestReceivingBandwidth(t *testing.T) {
	tap, eth, err := NewDevicesWithBandwidth(10000, 10000)
	if err != nil {
		t.Fatal(err)
	}
	defer tap.Close()

	frame := NewFrame(dest_mac, dev_mac)

	// Write the frame a lot of times so we get a more accurate idea of bandwidth.
	for i := 0; i < 100; i++ {
		err = eth.WriteFrame(frame)
		if err != nil {
			t.Fatal(err)
		}
	}

	start_time := time.Now()
	bytes_read := 0
	// Now read all those frames back.
	for i := 0; i < 100; i++ {
		read_this_cycle, err := waitForFrame(frame, tap)
		if err != nil {
			t.Fatal(err)
		}
		if read_this_cycle < 0 {
			t.Fatal("Reading frame failed.")
		}
		bytes_read += read_this_cycle
	}
	end_time := time.Now()

	// Check incoming bandwidth.
	elapsed := float64(end_time.Sub(start_time)) / float64(time.Second)
	bandwidth := float64(bytes_read) / elapsed

	if math.Abs(10000-bandwidth) > 50 {
		t.Fatalf("Got %f b/s of bandwidth, expected 10000.", bandwidth)
	}
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
	_, err = waitForFrame(frame, eth)
	if err != nil {
		t.Fatal(err)
	}
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
