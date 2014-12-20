package main

import (
	"encoding/hex"
	"log"
	"strings"
)

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
