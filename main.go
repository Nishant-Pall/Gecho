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

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			return
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid expected array length > 0")
			return
		}

		writer := NewWriter(conn)
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)

		writer.Write(result)
	}
}
