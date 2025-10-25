package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Nishant-Pall/Gecho/gloom"
)

var Handlers = map[string]func([]Value) Value{
	"PING":         pong,
	"GET":          get,
	"SET":          set,
	"HGET":         hget,
	"HSET":         hset,
	"HGETALL":      hgetall,
	"GLOOM_CREATE": gloomCreate,
	"GLOOM_ADD":    gloomAdd,
	"GLOOM_DELETE": gloomDelete,
	"GLOOM_LOOKUP": gloomLookup,
}

func pong([]Value) Value {
	return Value{typ: "string", str: "PONG"}
}

var SETsMu = sync.RWMutex{}
var SETs = map[string]string{}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `get` command"}
	}
	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: value}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `set` command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.RLock()
	SETs[key] = value
	SETsMu.RUnlock()

	return Value{typ: "string", str: "OK"}
}

var HSETSMu = sync.RWMutex{}
var HSETs = map[string]map[string]string{}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `hget` command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETSMu.Lock()
	value, ok := HSETs[hash][key]
	HSETSMu.Unlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `hset` command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETSMu.Lock()

	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETSMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `hgetall` command"}
	}

	hash := args[0].bulk

	HSETSMu.Lock()
	value, ok := HSETs[hash]
	HSETSMu.Unlock()

	if !ok {
		return Value{typ: "null"}
	}

	values := []Value{}

	for k, v := range value {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: values}
}

var gloomFilter *gloom.GloomFilter

func gloomCreate(args []Value) Value {

	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `gloom_create` command"}
	}

	len, err := strconv.Atoi(args[0].bulk)
	if err != nil {
		return Value{typ: "error", str: "Invalid input: length of gloom array"}
	}
	hashes, err := strconv.Atoi(args[1].bulk)
	if err != nil {
		return Value{typ: "error", str: "Invalid input: number of hashes"}
	}

	fmt.Printf("%v", gloomFilter)
	gloomFilter = gloom.NewGloomFilter()
	gloomFilter.CreateGloomFilter(uint64(len), uint64(hashes), gloom.MapHash)
	return Value{typ: "string", str: "OK"}
}

func gloomAdd(args []Value) Value {

	if gloomFilter == nil {
		return Value{typ: "error", str: "Gloom Filter not created"}
	}

	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `gloom_add` command"}
	}

	key := args[0].bulk

	err := gloomFilter.AddItem(key)

	if err != nil {
		return Value{typ: "error", str: fmt.Sprintf("%v", err)}
	}

	return Value{typ: "string", str: "ADDED"}
}

func gloomDelete(args []Value) Value {

	if gloomFilter == nil {
		return Value{typ: "error", str: "Gloom Filter not created"}
	}

	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `gloom_check` command"}
	}

	key := args[0].bulk

	err := gloomFilter.RemoveItem(key)

	if err != nil {
		return Value{typ: "error", str: fmt.Sprintf("%v", err)}
	}

	return Value{typ: "string", str: "Key removed"}
}
func gloomLookup(args []Value) Value {

	if gloomFilter == nil {
		return Value{typ: "error", str: "Gloom Filter not created"}
	}

	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `gloom_check` command"}
	}

	key := args[0].bulk

	ok, err := gloomFilter.Lookup(key)

	if err != nil {
		return Value{typ: "error", str: fmt.Sprintf("%v", err)}
	}

	return Value{typ: "string", str: strconv.FormatBool(ok)}
}
