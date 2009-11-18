// EBNF package for parsing IRC messages
//package irc_ebnf
package main

import (
    "bufio";
    "ebnf";
    "log";
    "os";
)

func main() {
    file, err := os.Open("irc.ebnf.go", os.O_RDONLY, 0666);
    if err != nil {
        log.Exit("Error: ", err)
    }
    reader := bufio.NewReader(file);
    src, err := reader.ReadBytes('\x00');
    if err != nil {
        log.Stdout("Error: ", err)
    }
    log.Stdout(string(src));
    // Read in the grammar
    grammar, err := ebnf.Parse("", src);
    log.Stdout(err);
    log.Stdout(grammar);
    // Verify the grammar
    err = ebnf.Verify(grammar, "message");
    log.Stdout(err);
}
