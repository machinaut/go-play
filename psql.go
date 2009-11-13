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
	hmesg[0] = 'F';
	binary.BigEndian.PutUint32(hmesg[1:5], uint32(3));
	binary.BigEndian.PutUint32(hmesg[5:9], uint32(len(str)));
	hmesg = bytes.Add(hmesg, mesg);
	fmt.Println(hmesg);
	fmt.Println(len(hmesg));
}
