package l2

import (
	"encoding/binary"
	"io"
	"net"
)

type SocketDevice struct {
	io.ReadWriteCloser
}

func NewListener(address string) (*SocketDevice, error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	c, err := ln.Accept()
	return &SocketDevice{c}, err
}

func NewDialer(address string) (*SocketDevice, error) {
	c, err := net.Dial("tcp", address)
	return &SocketDevice{c}, err
}

func (s *SocketDevice) WritePacket(data []byte) error {
	err := binary.Write(s, binary.BigEndian, int16(len(data)))
	if err != nil {
		return err
	}
	for written := 0; written < len(data); {
		n, err := s.Write(data[written:])
		if err != nil {
			return err
		}
		written += n
	}
	return nil
}

func (s *SocketDevice) ReadPacket() ([]byte, error) {
	var size int16
	err := binary.Read(s, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, size)
	_, err = io.ReadFull(s, buffer)
	return buffer, err
}
