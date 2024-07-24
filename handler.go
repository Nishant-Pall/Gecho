package main

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

var SETs = map[string]string{}
var SETmut sync.RWMutex

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "invalid number of arguments for `get` command"}
	}

	key := args[0].bulk

	SETmut.RLock()
	value, ok := SETs[key]
	SETmut.RUnlock()

	if !ok {
		return Value{typ: "error", str: "no value found for key"}
	}

	return Value{typ: "string", str: value}
}

func set(args []Value) Value {

	if len(args) != 2 {
		return Value{typ: "error", str: "invalid number of arguments for `set` command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETmut.Lock()
	SETs[key] = value
	SETmut.Unlock()

	return Value{typ: "string", str: "OK"}
}

func ping(args []Value) Value {
	if len(args) != 0 {
		return Value{typ: "string", str: fmt.Sprintf("PONG %s", args[0].bulk)}
	}
	return Value{typ: "string", str: "PONG"}
}
