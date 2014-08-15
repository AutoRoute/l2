package main

import (
	"code.google.com/p/tuntap"
	"log"
	"os/exec"
)

type TapDevice struct {
	dev *tuntap.Interface
}

func NewTapDevice() *TapDevice {
	fd, err := tuntap.Open("tap0", tuntap.DevTap)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Name:", fd.Name())

	ip_path, err := exec.LookPath("ip")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(ip_path, "link", "set", "dev", "tap0", "address", "00:24:d7:3e:71:b4")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(output))
		log.Fatal(err.Error())
	}

	cmd = exec.Command(ip_path, "link", "set", "dev", "tap0", "up")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Print(string(output))
		log.Fatal(err.Error())
	}
	return &TapDevice{fd}
}

func (t *TapDevice) Name() string {
	return t.dev.Name()
}

func (t *TapDevice) Close() {
	t.dev.Close()
}

func (t *TapDevice) ReadPacket() []byte {
	p, err := t.dev.ReadPacket()
	if err != nil {
		log.Fatal(err)
	}
	return p.Packet
}
