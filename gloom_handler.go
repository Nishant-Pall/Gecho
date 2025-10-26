package main

import (
	"fmt"
	"strconv"

	"github.com/Nishant-Pall/Gecho/gloom"
)

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

	gloomFilter = gloom.NewGloomFilter()
	gloomFilter.CreateGloomFilter(uint64(len), uint64(hashes), gloom.MapHash)
	return Value{typ: "string", str: "OK"}
}

func GloomLookup(args []Value) Value {

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

func GloomAdd(args []Value) Value {

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

func GloomDelete(args []Value) Value {

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
