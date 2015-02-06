package l2

import (
	"encoding/hex"
	"fmt"
	"log"
)

// Logs all frames which transit it.
type FrameLogger struct {
	d FrameReader
}

func (l FrameLogger) ReadFrame() ([]byte, error) {
	for {
		p, err := l.d.ReadFrame()
		if err == nil {
			PrintFrame(fmt.Sprint(l.d), p)
		} else {
			log.Print("Err reading frame:", err)
		}
		return p, err
	}
}

func (l FrameLogger) String() string {
	return "Logger{" + fmt.Sprint(l.d) + "}"
}

// Utility function to pretty print a ethernet frame plus a header
func PrintFrame(name string, data []byte) {
	E := EthFrame(data)
	log.Printf("%s: %s->%s %d bytes protocol %d",
		name,
		hex.EncodeToString(E.Source()),
		hex.EncodeToString(E.Destination()),
		len(data),
		E.Type())
}
