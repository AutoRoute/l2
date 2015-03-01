package l2

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

var (
	defaultdest  = MacToBytesOrDie("ff:ff:ff:ff:ff:ff")
	defaultframe = NewEthFrame(defaultdest, MacToBytesOrDie("00:00:00:00:00:00"), 1, make([]byte, 100))
	altframe     = NewEthFrame(
		MacToBytesOrDie("aa:bb:cc:dd:ee:00"), MacToBytesOrDie("00:00:00:00:00:00"), 1, make([]byte, 100))
)

type testReader []EthFrame

func (t *testReader) ReadFrame() (EthFrame, error) {
	if len(*t) > 0 {
		p := (*t)[0]
		if len(*t) > 1 {
			*t = (*t)[1:]
		} else {
			*t = nil
		}
		return p, nil
	}
	return nil, errors.New("Exhausted testing frames")
}

type readerTestCase struct {
	create    func(FrameReader) FrameReader
	input     []EthFrame
	output    []EthFrame
	stringrep []string
}

func createLogger(r FrameReader) FrameReader {
	return &FrameLogger{r}
}

func createFilter(r FrameReader) FrameReader {
	return NewFilterReader(r, defaultdest)
}

func TestReaders(t *testing.T) {
	testcases := []readerTestCase{
		{createLogger,
			[]EthFrame{defaultframe},
			[]EthFrame{defaultframe},
			[]string{"Logger"},
		},
		{createFilter,
			[]EthFrame{defaultframe, altframe},
			[]EthFrame{defaultframe},
			[]string{"filterReader", "ffffffffffff"},
		},
	}

	for _, tc := range testcases {
		tr := testReader(tc.input)
		reader := tc.create(&tr)

		// Check for all expected output
		for _, output := range tc.output {
			p, err := reader.ReadFrame()
			if err != nil {
				t.Errorf("Reader %v Expected %v error: %v", reader, output, err)
			}
			if !bytes.Equal(output, p) {
				t.Errorf("Reader %v Expected %v != %v", reader, output, p)
			}
		}

		// Once no input is left, testReader throws an error which a sane reader
		// should produce.
		_, err := reader.ReadFrame()
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
