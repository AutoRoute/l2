package main

import (
	"log"
	"net"
	"unsafe"
)

/*
#include <sys/socket.h>
#include <linux/if_packet.h>
#include <linux/if_ether.h>
#include <linux/if_arp.h>
#include <string.h>
#include <netinet/in.h>

socklen_t LLSize() {
    return sizeof(struct sockaddr_ll);
}

int SendData(int socket, int interface, char * data, size_t length) {
struct sockaddr_ll socket_address;
socket_address.sll_ifindex = interface;
socket_address.sll_halen = ETH_ALEN;
memcpy(socket_address.sll_addr, data, ETH_ALEN);
return sendto(socket, data, length, 0, (struct sockaddr*)&socket_address, LLSize());
};
*/
import "C"

type EthDevice struct {
	dev  C.int
	name string
}

func (e *EthDevice) Name() string {
	return e.name
}

func ConnectEthDevice(device string) *EthDevice {
	sock := C.socket(C.AF_PACKET, C.SOCK_RAW, C.int(C.htons(C.ETH_P_ALL)))
	ll_addr := C.struct_sockaddr_ll{}
	i, err := net.InterfaceByName(device)
	if err != nil {
		log.Fatal(err)
	}
	ll_addr.sll_ifindex = C.int(i.Index)
	C.bind(sock, (*C.struct_sockaddr)(unsafe.Pointer(&ll_addr)), C.LLSize())
	return &EthDevice{sock, device}
}

func (e *EthDevice) ReadPacket() []byte {
	buffer := [1523]byte{}
	n, err := C.recvfrom(e.dev, unsafe.Pointer(&buffer[0]), C.size_t(1523), 0, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	return buffer[0:n]
}
