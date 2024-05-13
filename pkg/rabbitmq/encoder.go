package rabbitmq

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

type Message map[string]interface{}

func Encode(m Message) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println("failed gob encode", err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func Decode(str string) Message {
	m := Message{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("failed base64 decode", err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		fmt.Println("failed gob decode", err)
	}
	return m
}
