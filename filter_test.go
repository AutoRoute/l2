package main

import (
	"bytes"
	"errors"
	"testing"
)

type testReader []EthPacket

func (t *testReader) ReadPacket() ([]byte, error) {
	if len(*t) > 0 {
		p := (*t)[0]
		if len(*t) > 1 {
			*t = (*t)[1:]
		} else {
			*t = nil
		}
		return p, nil
	}
	return nil, errors.New("Exhausted testing packets")
}

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
}
