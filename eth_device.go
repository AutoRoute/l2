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
	num  C.int
}

func (e *EthDevice) Name() string {
	return e.name
}

func ConnectEthDevice(device string) *EthDevice {
	sock, err := C.socket(C.AF_PACKET, C.SOCK_RAW, C.int(C.htons(C.ETH_P_ALL)))
	if err != nil {
		log.Fatal(err)
	}

	ll_addr := C.struct_sockaddr_ll{}
	i, err := net.InterfaceByName(device)
	if err != nil {
		log.Fatal(err)
	}
	ll_addr.sll_family = C.AF_PACKET
	ll_addr.sll_ifindex = C.int(i.Index)
	ll_addr.sll_protocol = C.__be16(C.htons(C.ETH_P_ALL))
	ll_addr.sll_pkttype = C.PACKET_HOST | C.PACKET_BROADCAST
	ok, err := C.bind(sock, (*C.struct_sockaddr)(unsafe.Pointer(&ll_addr)), C.LLSize())
	if ok != 0 || err != nil {
		log.Print("Error setting up eth device", device)
		log.Fatal(err)
	}
	return &EthDevice{sock, device, C.int(i.Index)}
}

func (e *EthDevice) ReadPacket() []byte {
	buffer := [1523]byte{}
	n, err := C.recvfrom(e.dev, unsafe.Pointer(&buffer[0]), C.size_t(1523), 0, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	return buffer[0:n]
}

func (e *EthDevice) WritePacket(data []byte) {
	socket_address := C.struct_sockaddr_ll{}
	socket_address.sll_ifindex = e.num
	socket_address.sll_halen = C.ETH_ALEN
	_, err := C.memcpy(unsafe.Pointer(&socket_address.sll_addr[0]),
		unsafe.Pointer(&data[0]), C.ETH_ALEN)
	if err != nil {
		log.Fatal(err)
	}
	ok, err := C.sendto(e.dev, unsafe.Pointer(&data[0]), C.size_t(len(data)),
		0, (*C.struct_sockaddr)(unsafe.Pointer(&socket_address)), C.LLSize())
	if ok != 0 || err != nil {
		log.Fatal(err)
	}
}
