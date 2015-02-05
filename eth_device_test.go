package l2

import (
	"testing"
)

func TestEthDevice(t *testing.T) {
	// Make sure if we open a non existent device it breaks
	_, err := ConnectEthDevice("Non existent device...")
	if err == nil {
		t.Fatal("Expected error to not be nil")
	}
}
