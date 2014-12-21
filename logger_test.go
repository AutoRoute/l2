package main

import (
	"bytes"
	"testing"
)

func TestLogger(t *testing.T) {
	dest := MacToBytesOrDie("ff:ff:ff:ff:ff:ff")
	packet := NewEthPacket(dest, MacToBytesOrDie("00:00:00:00:00:00"), 1, make([]byte, 100))

	// Make sure if we allow it through it comes through.
	tr := testReader([]EthPacket{packet})
	logger := PacketLogger{&tr}
	p, err := logger.ReadPacket()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p, packet) {
		t.Fatal("packet recieved is not what was sent")
	}
}
