package l2

import (
	"encoding/binary"
)

// A utility type to make introspecting ethernet frames easier.
type EthFrame []byte

func NewEthFrame(destination, source []byte, protocol uint16, data []byte) []byte {
	p := make([]byte, 12+2+len(data))
	copy(p[0:6], destination)
	copy(p[6:12], source)
	binary.BigEndian.PutUint16(p[12:14], protocol)
	copy(p[14:], data)
	return p
}

func (p EthFrame) Destination() []byte {
	return p[0:6]
}

func (p EthFrame) Source() []byte {
	return p[6:12]
}

// The ethernet type. Note that if this is <1504 it is likely a length instead
// and you are communicating with an extremely non standard ethernet device.
func (p EthFrame) Type() uint16 {
	return binary.BigEndian.Uint16(p[12:14])
}

func (p EthFrame) Data() []byte {
	return p[14:]
}
