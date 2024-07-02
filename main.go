package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func setupCon() {
	fmt.Println("Initiating connection...")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening at: ", l.Addr().String())

	conn, _err := l.Accept()
	if _err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		buff := make([]byte, 1024)

		_, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}
		conn.Write([]byte("+OK\r\n"))
	}

}

func main() {
	input := "$7\r\nNishant\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	b, _ := reader.ReadByte()

	if b != '$' {
		fmt.Println("Invalid type, expecting bulk strings only")
		os.Exit(1)
	}

	size, _ := reader.ReadByte()
	sizeStr, _ := strconv.ParseInt(string(size), 10, 64)
	fmt.Println(sizeStr)

	reader.ReadByte()
	reader.ReadByte()

	name := make([]byte, sizeStr)
	reader.Read(name)

	fmt.Println(string(name))
}
