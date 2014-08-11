package main

import (
    "code.google.com/p/tuntap"
    "log"
    "fmt"
    "os/exec"
)

// #include <sys/socket.h>
// #include <linux/if_packet.h>
// #include <linux/if_ether.h>
// #include <linux/if_arp.h>
// #include <string.h>
//
// int CreateSocket() {
//   return socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
// };
//
// int SendData(int socket, int interface, char * data, size_t length) {
//   struct sockaddr_ll socket_address;
//   socket_address.sll_ifindex = interface;
//   socket_address.sll_halen = ETH_ALEN;
//   memcpy(socket_address.sll_addr, data, ETH_ALEN);
//   return sendto(socket, data, length, 0, (struct sockaddr*)&socket_address, sizeof(socket_address));
// };
//
// int ReceiveData(int socket, void* buffer) {
//   return recvfrom(socket, buffer, ETH_FRAME_LEN, 0, NULL, NULL);
// }
import "C"

func print_recv(tap *tuntap.Interface) {
    for {
        p, err := tap.ReadPacket()
        if err != nil {
            log.Fatal(err)
        }
        log.Print(p)
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
    fmt.Scanln();
}
