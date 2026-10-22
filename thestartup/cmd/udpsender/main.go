package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.DialUDP(addr.Network(), nil, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		l, err := reader.ReadString(0x000A)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = conn.Write([]byte(l))
		if err != nil {
			fmt.Println(err)
			continue
		}

	}
}
