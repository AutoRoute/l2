package l2

import (
	"code.google.com/p/tuntap"
	"log"
	"os/exec"
)

// A Tap Device is a new device that this program has created. In this case
// the normal semantics are inverted, in that frames sent to the device
// are what this interface will read and vice versa.
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

func (t *TapDevice) String() string {
	return "TapDevice{" + t.dev.Name() + "}"
}

func (t *TapDevice) Close() {
	t.dev.Close()
}

func (t *TapDevice) ReadFrame() ([]byte, error) {
	p, err := t.dev.ReadPacket()
	if err != nil {
		return nil, err
	}
	return p.Packet, nil
}

func (t *TapDevice) WriteFrame(data []byte) error {
	t.dev.WritePacket(
		&tuntap.Packet{
			Protocol: int(EthFrame(data).Type()),
			Packet:   data})
	return nil
}
