package main

type PacketReader interface {
	ReadPacket() []byte
}

type PacketWriter interface {
	WritePacket([]byte)
}

type PacketDevice interface {
	PacketReader
	Device
}

type Device interface {
	Name() string
}
