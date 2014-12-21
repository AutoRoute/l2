package main

import (
	"errors"
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
