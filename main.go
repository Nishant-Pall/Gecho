package main

import (
	"fmt"
	"net"
	"strings"
)

const aofPath = "dump.aof"

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Printf("cannot initiate listener: %v", err)
		return
	}

	fmt.Printf("Initiated %v connection at address %v \n", listener.Addr().Network(), listener.Addr().String())

	aof, err := NewAof(aofPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	defer aof.Close()

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("Cannot initiate connection: %v", err)
		return
	}

	defer conn.Close()

	for {
		resp := NewReader(conn)

		value, err := resp.Read()

		fmt.Println(value)
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			return
		}

		if len(value.array) == 0 {
			fmt.Println("Unexpected length")
			return
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)
		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		fmt.Println(result)
		writer.Write(result)
	}
}
