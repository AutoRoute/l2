package main

type PacketReader interface {
	ReadPacket() []byte
}

type PacketWriter interface {
	WritePacket([]byte)
}

type PacketDevice interface {
	PacketReader
	PacketWriter
	Name() string
}
