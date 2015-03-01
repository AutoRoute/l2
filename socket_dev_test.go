package l2

import (
	"net"
	"testing"
)

func newListener(address string) (FrameReadWriter, error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	c, err := ln.Accept()
	return WrapReadWriter(c), err
}

func newDialer(address string) (FrameReadWriter, error) {
	c, err := net.Dial("tcp", address)
	return WrapReadWriter(c), err
}

func TestSocketDevice(t *testing.T) {
	sync := make(chan struct{})
	go func() {
		sync <- struct{}{}
		c, err := newListener("127.0.0.1:9000")
		if err != nil {
			t.Fatal(err)
		}
		p, err := c.ReadFrame()
		if err != nil {
			t.Fatal(err)
		}
		err = c.WriteFrame(p)
		if err != nil {
			t.Fatal(err)
		}
	}()
	<-sync
	c, err := newDialer("127.0.0.1:9000")
	if err != nil {
		t.Fatal(err)
	}
	x := "1234567891abcdefgh"
	err = c.WriteFrame(EthFrame([]byte(x)))
	if err != nil {
		t.Fatal(err)
	}
	p, err := c.ReadFrame()
	if string(p) != x {
		t.Fatalf("Expected %s got %s", x, string(p))
	}
}
