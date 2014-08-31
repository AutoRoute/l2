package main

import (
	"code.google.com/p/tuntap"
	"log"
	"os/exec"
)

type TapDevice struct {
	dev *tuntap.Interface
}

func NewTapDevice(mac, dev string) (*TapDevice, error) {
	fd, err := tuntap.Open(dev, tuntap.DevTap)
	if err != nil {
        return nil, err
	}

	ip_path, err := exec.LookPath("ip")
	if err != nil {
        return nil, err
	}

	cmd := exec.Command(ip_path, "link", "set", "dev", dev, "address", mac)
	output, err := cmd.CombinedOutput()
	if err != nil {
        log.Print("Command output:", string(output))
        return nil, err
	}

	cmd = exec.Command(ip_path, "link", "set", "dev", dev, "up")
	output, err = cmd.CombinedOutput()
	if err != nil {
        log.Print("Command output:", string(output))
        return nil, err
	}
	return &TapDevice{fd}, nil
}

func (t *TapDevice) Name() string {
	return t.dev.Name()
}

func (t *TapDevice) Close() {
	t.dev.Close()
}

func (t *TapDevice) ReadPacket() ([]byte, error) {
	p, err := t.dev.ReadPacket()
	if err != nil {
        return nil, err
	}
	return p.Packet, nil
}

func (t *TapDevice) WritePacket(data []byte) error {
	t.dev.WritePacket(
		&tuntap.Packet{
			Protocol: EthPacket(data).TypeInt(),
			Packet:   data})
    return nil
}
