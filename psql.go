package main

import (
	"fmt";
	"log";
	"net";
	"strings";
	"bytes";
	"encoding/binary";
)

func main() {
	_, err := net.Dial("tcp", "", "152.7.16.113:5432");
	if err != nil {
		log.Exitf("Error while connecting: %s", err)
	}
	log.Stdout("I think we connected");
	fmt.Println(len("hello\x00"));
	str := "user\x00alex\x00database\x00ar_test\x00\x00";
	mesg := strings.Bytes(str);
	hmesg := make([]byte, 1+4+4);
	binary.BigEndian.PutUint32(hmesg[1:5], uint32(3));
	fmt.Println(str, mesg, hmesg);
	fmt.Println(len(str), len(mesg), len(hmesg));
    bytes.Add
}
