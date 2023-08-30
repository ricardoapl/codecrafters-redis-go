package main

import (
	"bufio"
	"errors"
	"strconv"
)

const (
	Array        byte = '*'
	BulkString   byte = '$'
	SimpleString byte = '+'
)

type RESP struct {
	Type     byte
	Value    []byte
	Elements []RESP
}

func Deserialize(rd *bufio.Reader) (RESP, error) {
	t, err := rd.ReadByte()
	if err != nil {
		return RESP{}, err
	}
	switch t {
	case Array:
		return deserializeArray(rd)
	case BulkString:
		return deserializeBulkString(rd)
	case SimpleString:
		return deserializeSimpleString(rd)
	}
	return RESP{}, errors.New("unknown RESP data type")
}

func deserializeArray(rd *bufio.Reader) (RESP, error) {
	v, _, err := rd.ReadLine()
	if err != nil {
		return RESP{}, err
	}
	count, err := strconv.Atoi(string(v))
	if err != nil {
		return RESP{}, err
	}
	if count == 0 {
		return RESP{Type: Array}, nil
	}
	elements := []RESP{}
	for c := 0; c < count; c++ {
		r, err := Deserialize(rd)
		if err != nil {
			return RESP{}, err
		}
		elements = append(elements, r)
	}
	return RESP{
		Type:     Array,
		Elements: elements,
	}, nil
}

func deserializeBulkString(rd *bufio.Reader) (RESP, error) {
	v, _, err := rd.ReadLine()
	if err != nil {
		return RESP{}, err
	}
	length, err := strconv.Atoi(string(v))
	if err != nil {
		return RESP{}, err
	}
	if length == 0 {
		return RESP{Type: BulkString}, nil
	}
	v, _, err = rd.ReadLine()
	if err != nil {
		return RESP{}, err
	}
	return RESP{
		Type:  BulkString,
		Value: v,
	}, nil
}

func deserializeSimpleString(rd *bufio.Reader) (RESP, error) {
	v, _, err := rd.ReadLine()
	if err != nil {
		return RESP{}, err
	}
	return RESP{
		Type:  SimpleString,
		Value: v,
	}, nil
}
