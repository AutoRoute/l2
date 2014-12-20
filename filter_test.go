package main

import (
	"bytes"
	"testing"
)

type testReader EthPacket

func (t testReader) ReadPacket() ([]byte, error) {
	return EthPacket(t), nil
}

func TestFilter(t *testing.T) {
	dest := MacToBytesOrDie("ff:ff:ff:ff:ff:ff")
	packet := NewEthPacket(dest, MacToBytesOrDie("00:00:00:00:00:00"), 1, make([]byte, 100))
	filter := NewFilterPacket(testReader(packet), dest)
	p, err := filter.ReadPacket()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p, packet) {
		t.Fatal("packet recieved is not what was sent")
	}
}
