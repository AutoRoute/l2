package main

import (
    "code.google.com/p/tuntap"
    "encoding/hex"
    "fmt"
    "log"
    "net"
    "os/exec"
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

func PrintPacket(name string, data []byte) {
    log.Printf("%s: %s->%s", name,
        hex.EncodeToString(data[0:6]),
        hex.EncodeToString(data[6:12]))
}
func print_recv(tap *tuntap.Interface) {
    for {
        p, err := tap.ReadPacket()
        if err != nil {
            log.Fatal(err)
        }
        PrintPacket("tap0", p.Packet)
    }
}

func CreateListenSocket(listen_interface string) C.int {
    sock := C.socket(C.AF_PACKET, C.SOCK_RAW, C.int(C.htons(C.ETH_P_ALL)))
    ll_addr := C.struct_sockaddr_ll{}
    i, err := net.InterfaceByName(listen_interface)
    if err != nil {
        log.Fatal(err)
    }
    ll_addr.sll_ifindex = C.int(i.Index)
    C.bind(sock, (*C.struct_sockaddr)(unsafe.Pointer(&ll_addr)), C.LLSize())
    return sock
}

func cprint_recv() {
    sock := CreateListenSocket("wlp3s0")
    for {
        buffer := [1523]byte{}
        n, err := C.recvfrom(sock, unsafe.Pointer(&buffer[0]), C.size_t(1523), 0, nil, nil)
        if err != nil {
            log.Fatal(err)
        }
        PrintPacket("wlp3s0", buffer[0:n])
    }
}

func main() {
    fd, err := tuntap.Open("tap0", tuntap.DevTap)
    if err != nil {
        log.Fatal(err)
    }
    log.Print("Name:", fd.Name())
    defer fd.Close()

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


    go print_recv(fd)
    go cprint_recv()
    fmt.Scanln();
}
