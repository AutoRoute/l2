package l2

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// filterReader is a FrameReader which only allows through frames which match the list of
// frames is is supplied with.
type filterReader struct {
	mac    [][]byte
	device FrameReader
}

// Construct a filter which only allows through the specified mac addresses
func NewFilterFrame(dev FrameReader, mac ...[]byte) FrameReader {
	return filterReader{mac, dev}
}

func (f filterReader) ReadFrame() (EthFrame, error) {
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

func (f filterReader) String() string {
	s := "filterReader{" + fmt.Sprint(f.device)
	for _, mac := range f.mac {
		s += ", " + hex.EncodeToString(mac)
	}
	s += "}"
	return s
}
