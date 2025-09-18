package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": pong,
	"GET":  get,
	"SET":  set,
}

func pong([]Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func get([]Value) Value {
	mu := sync.Mutex

}
func set([]Value) Value {
	mu := sync.Mutex

}
