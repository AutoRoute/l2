package l2

import (
	"testing"
)

func TestSocketDevice(t *testing.T) {
	sync := make(chan struct{})
	go func() {
		sync <- struct{}{}
		c, err := NewListener("127.0.0.1:9000")
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
	c, err := NewDialer("127.0.0.1:9000")
	if err != nil {
		t.Fatal(err)
	}
	x := "1234567891abcdefgh"
	err = c.WriteFrame([]byte(x))
	if err != nil {
		t.Fatal(err)
	}
	p, err := c.ReadFrame()
	if string(p) != x {
		t.Fatalf("Expected %s got %s", x, string(p))
	}
}
