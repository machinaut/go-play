package main

import (
    "fmt";
    "io";
    "log";
    //"unicode";
    "utf8";
)

type rune int

func main() {
    // Read file into []byte
    bytes, err := io.ReadFile("math.peg");
    if err != nil {
        log.Exit("Error :")
    }

    // Convert []byte into []rune
    bmax := len(bytes);
    rmax := utf8.RuneCount(bytes); // # runes
    runes := make([]rune, rmax);
    for bi, ri := 0, 0; ri < rmax; ri++ {
        r, size := utf8.DecodeRune(bytes[bi:bmax]);
        bi += size;          // Advance in bytestream
        runes[ri] = rune(r); // Add rune
    }

    fmt.Println(runes);
}
