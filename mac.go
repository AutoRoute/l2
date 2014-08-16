package main

import (
	"encoding/hex"
	"log"
	"strings"
)

type MAC string

func (m MAC) ToBytes() []byte {
	r, err := hex.DecodeString(strings.Replace(string(m), ":", "", -1))
	if err != nil {
		log.Fatal(err)
	}
	return r
}
