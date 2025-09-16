package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Printf("cannot initiate listener: %v", err)
		return
	}

	fmt.Printf("Initiated %v connection at address %v \n", listener.Addr().Network(), listener.Addr().String())

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("Cannot initiate connection: %v", err)
		return
	}

	defer conn.Close()

	for {
		resp := NewReader(conn)

		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v \r\n", value)

		conn.Write([]byte("+OK\r\n"))
	}
}
