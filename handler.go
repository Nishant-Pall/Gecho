package main

import (
	"fmt"
	"strconv"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"INCR":    increment,
	"DCR":     decrement,
}

var SETs = map[string]string{}
var SETmut sync.RWMutex

var HSETs = map[string]map[string]string{}
var HSETmut sync.RWMutex

func decrement(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "invalid number of arguments for `incr` command"}
	}

	key := args[0].bulk

	if _, ok := SETs[key]; !ok {
		return Value{typ: "error", str: "Key does not exist"}
	}

	counter, _ := strconv.Atoi(SETs[key])
	counter = counter - 1

	SETs[key] = strconv.Itoa(counter)

	return Value{typ: "string", str: "OK"}
}

func increment(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "invalid number of arguments for `incr` command"}
	}

	key := args[0].bulk

	if _, ok := SETs[key]; !ok {
		return Value{typ: "error", str: "Key does not exist"}
	}

	counter, _ := strconv.Atoi(SETs[key])
	counter = counter + 1

	SETs[key] = strconv.Itoa(counter)

	return Value{typ: "string", str: "OK"}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "invalid number of arguments for `hset` command"}
	}

	fmt.Println(args)

	hash := args[0].bulk
	key2 := args[1].bulk
	value := args[2].bulk

	HSETmut.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key2] = value
	HSETmut.Unlock()

	return Value{typ: "string", str: "OK"}
}
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "invalid number of arguments for `hget` command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETmut.RLock()
	value, ok := HSETs[hash][key]
	HSETmut.RUnlock()

	if !ok {
		return Value{typ: "error", str: "no value found for key"}
	}

	return Value{typ: "string", str: value}
}
func hgetall(args []Value) Value {
	vmap := HSETs
	return Value{typ: "string", str: fmt.Sprint(vmap)}
}

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
