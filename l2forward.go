// Program l2forward is a simple binary to foward networking devies at the layer two level.
package l2

import (
	"encoding/hex"
	"log"
	"strings"
)

// One of the core interface abstraction in l2forward. This represents
// something which ethernet frames can be read from. This is distinct from
// io.Reader because you cannot slice l2 ethernet frames arbitrarily.
type FrameReader interface {
	ReadFrame() ([]byte, error)
}

// One of the core interface abstraction in l2forward. This represents
// something which ethernet frames can be written to. This is distinct from
// io.Reader because you cannot slice l2 ethernet frames arbitrarily.
type FrameWriter interface {
	WriteFrame([]byte) error
}

// Basic utility function to take a string and turn it into a mac address. Will
// die if the string is not valid.
func MacToBytesOrDie(m string) []byte {
	b, err := hex.DecodeString(strings.Replace(m, ":", "", -1))
	if err != nil {
		log.Fatal(err)
	}
	if len(b) != 6 {
		log.Fatal("Expected mac of length 6 bytes got ", len(b))
	}
	return b
}

// Local equivalent of io.Copy, will shove frames from a FrameReader
// into a FrameWriter
func SendFrames(source FrameReader, destination FrameWriter) {
	for {
		p, err := source.ReadFrame()
		if err != nil {
			log.Fatal("Failure to read from source", source, err)
		}
		if destination.WriteFrame(p) != nil {
			log.Fatal("Failure to write to", destination, err)
		}
	}
}
