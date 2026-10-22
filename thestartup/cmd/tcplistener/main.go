package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
	"os"
)

func main() {
	socket, err := net.Listen("tcp", ":42069")
	if err != nil {
		os.Exit(1)
	}
	defer socket.Close()
	for {
		conn, err := socket.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Connection accepted")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			continue
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range r.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(r.Body))
	}

}
