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
	pg_adr, err := net.ResolveTCPAddr("152.7.16.113:5432");
	if err != nil {
		log.Exitf("Error while resolving address: %s", err)
	}
	conn, err := net.DialTCP("tcp", nil, pg_adr);
	if err != nil {
		log.Exitf("Error while connecting: %s", err)
	}
	log.Stdout("I think we connected");

	strs := make([]string, 4);
	strs[0] = "user";
	strs[1] = "alex";
	strs[2] = "database";
	strs[3] = "ar_test";
	str2 := strings.Join(strs, "\x00") + "\x00\x00";

	mesg2 := strings.Bytes(str2);
	hmesg := make([]byte, 1+4+4);
	hmesg[0] = 'F';
	binary.BigEndian.PutUint32(hmesg[1:5], uint32(3));
	hmesg = bytes.Add(hmesg, mesg2);
	binary.BigEndian.PutUint32(hmesg[5:9], uint32(len(hmesg)-1));
	fmt.Println(hmesg);
	fmt.Println(len(hmesg));
	n, err := conn.Write(hmesg);

	if err != nil {
		log.Exitf("Error writing TCP: %s", err)
	}
	log.Stdoutf("wrote %d", n);

	result := make([]byte, 100);
	err = conn.SetReadTimeout(0);
	n, err = conn.Read(result);
	if err != nil {
		log.Stdoutf("Error reading TCP: %s", err)
	}
	log.Stdoutf("Read %d", n);
}
