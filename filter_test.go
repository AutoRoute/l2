package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	dest := MacToBytesOrDie("ff:ff:ff:ff:ff:ff")
	packet := NewEthPacket(dest, MacToBytesOrDie("00:00:00:00:00:00"), 1, make([]byte, 100))

	// Make sure if we allow it through it comes through.
	tr := testReader([]EthPacket{packet})
	filter := NewFilterPacket(&tr, dest)
	p, err := filter.ReadPacket()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p, packet) {
		t.Fatal("packet recieved is not what was sent")
	}
	p, err = filter.ReadPacket()
	if err == nil {
		t.Fatal("Expected error")
	}

	// Make sure if we don't, it doesn't.
	alt := MacToBytesOrDie("00:00:00:00:00:00")
	tr = testReader([]EthPacket{packet})
	filter = NewFilterPacket(&tr, alt)
	p, err = filter.ReadPacket()
	if err == nil {
		t.Fatal("Expected error")
	}

	// Make sure that the printing works correctly
	if !strings.Contains(filter.String(), "000000000000") {
		t.Fatal("Expected to see address in string rep", filter.String())
	}
	if !strings.Contains(filter.String(), "FilterPacket") {
		t.Fatal("Expected to see name in string rep", filter.String())
	}
}
