package main

import (
    "./psql_constants";
    "fmt";
    "log";
    "net";
    "strings";
    "bytes";
    "encoding/binary";
    "crypto/md5";
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
    hmesg := make([]byte, 4+4);
    binary.BigEndian.PutUint32(hmesg[4:8], uint32(3<<16));
    hmesg = bytes.Add(hmesg, mesg2);
    binary.BigEndian.PutUint32(hmesg[0:4], uint32(len(hmesg)));
    n, err := conn.Write(hmesg);

    if err != nil {
        log.Exitf("Error writing TCP: %s", err)
    }
    log.Stdoutf("wrote %d", n);

    result := make([]byte, 12);
    // the largest response we can get is 12 bytes
    if n, err = conn.Read(result); err != nil {
        log.Exitf("Error reading TCP (Read %d bytes): %s", n, err)
    }
    log.Stdoutf("%c", result[0]);

    switch string(result[0]) {
    case psql_constants.Authentication:
        log.Stdout("Got authentication message")
    default:
        log.Exit("Did not get authentication message")
    }

    mesglen := binary.BigEndian.Uint32(result[1:5]);
    mesgtype := binary.BigEndian.Uint32(result[5:9]);

    fmt.Println(mesglen, mesgtype);

    passhash := md5.New();
    passhash.Write(strings.Bytes("jar1!" + "alex"));
    salthash := md5.New();
    salthash.Write(passhash.Sum());
    salthash.Write(result[9:12]);

    passresponse := make([]byte, 5);
    passresponse[0] = 'p';
    binary.BigEndian.PutUint32(passresponse[1:5], uint32(4+salthash.Size()+1));
    passresponse = bytes.Add(passresponse, strings.Bytes("md5"));
    passresponse = bytes.Add(passresponse, salthash.Sum());
    passresponse = bytes.AddByte(passresponse, byte(0));
    fmt.Println(passresponse);
    n, err = conn.Write(passresponse);
    if err != nil {
        log.Exitf("Error writing TCP: %s", err)
    }
    log.Stdoutf("wrote %d", n);

    result = make([]byte, 18);
    // the largest response we can get is 12 bytes
    if n, err = conn.Read(result); err != nil {
        log.Exitf("Error reading TCP (Read %d bytes): %s", n, err)
    }
    fmt.Println(result);

    conn.Close();
}
