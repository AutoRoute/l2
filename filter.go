package l2

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// FilterFrame is a FrameReader which only allows through frames which match the list of
// frames is is supplied with.
type FilterFrame struct {
	mac    [][]byte
	device FrameReader
}

// Construct a filter which only allows through the specified mac addresses
func NewFilterFrame(dev FrameReader, mac ...[]byte) *FilterFrame {
	return &FilterFrame{mac, dev}
}

func (f FilterFrame) ReadFrame() (EthFrame, error) {
	for {
		p, err := f.device.ReadFrame()
		if err != nil {
			return p, err
		}
		for _, mac := range f.mac {
			if bytes.Equal(EthFrame(p).Destination(), mac) {
				return p, nil
			}
		}
	}
}

func (f FilterFrame) String() string {
	s := "FilterFrame{" + fmt.Sprint(f.device)
	for _, mac := range f.mac {
		s += ", " + hex.EncodeToString(mac)
	}
	s += "}"
	return s
}
