package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Initiating connection...")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid Command: ", command)
			return
		}
		handler(args)
	})

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

		if value.typ != "Array" {
			fmt.Println("Invalid request, expected array")
			return
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid expected array length > 0")
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

		writer.Write(result)
	}
}
