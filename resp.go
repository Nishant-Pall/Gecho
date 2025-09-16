package main

import (
	"bufio"
	"io"
	"strconv"
)

const (
	STRING          = '+'
	ERROR           = '-'
	INTEGER         = ':'
	BULK            = '$'
	ARRAY           = '*'
	NULLS           = '_'
	BOOLEAN         = '#'
	DOUBLE          = ','
	BIG_NUMBER      = '('
	BULK_ERROR      = '!'
	VERBATIM_STRING = '='
	MAP             = '%'
	ATTRIBUTE       = '|'
	SET             = '~'
	PUSH            = '>'
)

type Resp struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		typ, err := r.reader.ReadByte()
		if err != nil {
			return nil, n, err
		}

		n += 1
		line = append(line, typ)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()

	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)

	if err != nil {
		return 0, n, err
	}

	return int(i64), n, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)
	r.reader.Read(bulk)
	v.bulk = string(bulk)

	r.readLine()

	return v, nil
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	v.array = make([]Value, length)

	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val)
	}

	return v, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		return Value{}, nil
	}
}

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

func (v *Value) MarshallString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v *Value) MarshallBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v *Value) MarshallArray() []byte {
	var bytes []byte
	length := len(v.array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(length)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < length; i++ {
		bytes = append(bytes, v.array[i].Marshall()...)
	}

	return bytes
}

func (v *Value) MarshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v *Value) MarshallNull() []byte {

	return []byte("$-1\r\n")
}

func (v *Value) Marshall() []byte {
	switch v.typ {
	case "array":
		return v.MarshallArray()
	case "bulk":
		return v.MarshallBulk()
	case "string":
		return v.MarshallString()
	case "null":
		return v.MarshallNull()
	case "error":
		return v.MarshallError()
	default:
		return []byte{}

	}
}
