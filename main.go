package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Initiating connection...")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening at: ", l.Addr().String())
	conn, _err := l.Accept()
	if _err != nil {
		fmt.Println(_err)
		return
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			fmt.Println("error reading from client: ", err.Error())
			return
		}

		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))
	}
}
