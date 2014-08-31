package main

type PacketReader interface {
	ReadPacket() ([]byte, error)
}

type PacketWriter interface {
	WritePacket([]byte) error
}

type Namer interface {
	Name() string
}

type PacketDevice interface {
	PacketReader
	PacketWriter
	Namer
}

type PacketReaderDevice interface {
	PacketReader
	Namer
}

type PacketWriterDevice interface {
	PacketWriter
	Namer
}
