package main

import (
    "./psql_constants";
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

    strs := make([]string, 6);
    strs[0] = "user";
    strs[1] = "alex";
    strs[2] = "password";
    strs[2] = "jar1!";
    strs[4] = "database";
    strs[5] = "ar_test";
    str2 := strings.Join(strs, "\x00") + "\x00\x00";

    mesg2 := strings.Bytes(str2);
    hmesg := make([]byte, 4+4);
    binary.BigEndian.PutUint32(hmesg[4:8], uint32(3<<16));
    hmesg = bytes.Add(hmesg, mesg2);
    binary.BigEndian.PutUint32(hmesg[0:4], uint32(len(hmesg)));
    fmt.Println(hmesg);
    fmt.Println(len(hmesg));
    n, err := conn.Write(hmesg);

    if err != nil {
        log.Exitf("Error writing TCP: %s", err)
    }
    log.Stdoutf("wrote %d", n);

    result := make([]byte, 12);
    // the largest response we can get is 12 bytes
    n, err = conn.Read(result);
    if err != nil {
        log.Stdoutf("Error reading TCP (Read %d bytes): %s", n, err)
    }
    fmt.Println(result);

    fmt.Println(psql_constants.Authentication);

    conn.Close();

}
