package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

var (
	defaultdest   = MacToBytesOrDie("ff:ff:ff:ff:ff:ff")
	defaultpacket = NewEthPacket(defaultdest, MacToBytesOrDie("00:00:00:00:00:00"), 1, make([]byte, 100))
	altpacket     = NewEthPacket(
		MacToBytesOrDie("aa:bb:cc:dd:ee:00"), MacToBytesOrDie("00:00:00:00:00:00"), 1, make([]byte, 100))
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

type readerTestCase struct {
	create    func(PacketReader) PacketReader
	input     []EthPacket
	output    []EthPacket
	stringrep []string
}

func createLogger(r PacketReader) PacketReader {
	return &PacketLogger{r}
}

func createFilter(r PacketReader) PacketReader {
	return NewFilterPacket(r, defaultdest)
}

func TestReaders(t *testing.T) {
	testcases := []readerTestCase{
		{createLogger,
			[]EthPacket{defaultpacket},
			[]EthPacket{defaultpacket},
			[]string{"Logger"},
		},
		{createFilter,
			[]EthPacket{defaultpacket, altpacket},
			[]EthPacket{defaultpacket},
			[]string{"FilterPacket", "ffffffffffff"},
		},
	}

	for _, tc := range testcases {
		tr := testReader(tc.input)
		reader := tc.create(&tr)

		// Check for all expected output
		for _, output := range tc.output {
			p, err := reader.ReadPacket()
			if err != nil {
				t.Errorf("Reader %v Expected %v error: %v", reader, output, err)
			}
			if !bytes.Equal(output, p) {
				t.Errorf("Reader %v Expected %v != %v", reader, output, p)
			}
		}

		// Once no input is left, testReader throws an error which a sane reader
		// should produce.
		_, err := reader.ReadPacket()
		if err == nil {
			t.Errorf("Reader %v Expected error got %v", reader, err)
		}

		// Make sure the string represenation is sane
		for _, r := range tc.stringrep {
			o := fmt.Sprint(reader)
			if !strings.Contains(o, r) {
				t.Errorf("Reader %v should contain %s", reader, o)
			}
		}
	}
}
