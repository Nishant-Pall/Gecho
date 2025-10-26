package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Nishant-Pall/Gecho/gloom"
)

var GloomMU = sync.RWMutex{}
var gloomFilter *gloom.BaseGloomFilter

func GloomCreate(args []Value) Value {

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

	GloomMU.Lock()
	filter, err := gloom.CreateGloomFilter(uint64(len), uint64(hashes), gloom.MapHash)
	gloomFilter = filter
	GloomMU.Unlock()

	if err != nil {
		return Value{typ: "error", str: "Error creating Gloom Filter"}
	}

	return Value{typ: "string", str: fmt.Sprintf("Gloom filter created of length: %v", gloomFilter.Len())}
}

func GloomLookup(args []Value) Value {

	if gloomFilter == nil {
		return Value{typ: "error", str: "Gloom Filter not created"}
	}

	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `gloom_check` command"}
	}

	key := args[0].bulk

	GloomMU.Lock()
	ok, err := gloomFilter.Lookup(key)
	GloomMU.Unlock()

	if err != nil {
		return Value{typ: "error", str: fmt.Sprintf("%v", err)}
	}

	return Value{typ: "string", str: strconv.FormatBool(ok)}
}

func GloomAdd(args []Value) Value {

	if gloomFilter == nil {
		return Value{typ: "error", str: "Gloom Filter not created"}
	}

	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `gloom_add` command"}
	}

	key := args[0].bulk

	GloomMU.Lock()
	err := gloomFilter.AddItem(key)
	GloomMU.Unlock()

	if err != nil {
		return Value{typ: "error", str: fmt.Sprintf("%v", err)}
	}

	return Value{typ: "string", str: "ADDED"}
}

func GloomDelete(args []Value) Value {

	if gloomFilter == nil {
		return Value{typ: "error", str: "Gloom Filter not created"}
	}

	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for `gloom_check` command"}
	}

	key := args[0].bulk

	GloomMU.Lock()
	err := gloomFilter.RemoveItem(key)
	GloomMU.Unlock()

	if err != nil {
		return Value{typ: "error", str: fmt.Sprintf("%v", err)}
	}

	return Value{typ: "string", str: "Key removed"}
}
