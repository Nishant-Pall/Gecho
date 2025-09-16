package main

var Handlers = map[string]func([]Value) Value{
	"PING": pong,
}

func pong([]Value) Value {
	return Value{typ: "string", str: "PONG"}
}
