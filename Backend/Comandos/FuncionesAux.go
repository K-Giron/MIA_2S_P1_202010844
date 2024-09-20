package Comandos

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Codifica de Struct a []Bytes
func struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)

	// ERROR
	if err != nil && err != io.EOF {
		MsgError(err)
	}

	return buf.Bytes()
}

func LeerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		MsgError(err)
	}
	return bytes
}

func EscribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

// Decodifica de [] Bytes a Struct
func bytes_a_struct_mbr(s []byte) Mbr {
	p := Mbr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		MsgError(err)
	}

	return p
}

func bytes_a_struct_ebr(s []byte) Ebr {
	// Descodificacion
	p := Ebr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		MsgError(err)
	}

	return p
}

// Muestra el mensaje de error
func MsgError(err error) {
	fmt.Println("ERROR: ", err)
}

func Error(op string, mensaje string) {
	fmt.Println("ERROR: ", op, mensaje)
}

func Comparar(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}
