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
    file, err := os.Open("full_irc.ebnf.go", os.O_RDONLY, 0666);
    if err != nil {
        log.Exit("Error: ", err)
    }
    reader := bufio.NewReader(file);
    src, err := reader.ReadBytes('\x00');
    if err != nil {
        log.Stdout("Reading stopped: ", err)
    }
    // Read in the grammar
    //log.Stdoutf("%s", src);
    grammar, err := ebnf.Parse("", src);
    if err != nil {
        log.Exit("Parse Error: ", err)
    }
    // Verify the grammar
    err = ebnf.Verify(grammar, "message");
    if err != nil {
        log.Exit("Verification Error: ", err)
    }
    log.Stdout("Success");
}
