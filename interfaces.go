package main

type PacketReader interface {
	ReadPacket() ([]byte, error)
}

type PacketWriter interface {
	WritePacket([]byte) error
}
