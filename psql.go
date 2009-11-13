package main

import (
	"fmt";
	"net";
)

func main() {
	_, e := net.Dial("tcp", "", "152.7.16.113:5432");
	if e != nil {
		fmt.Println("We got ourselves a problem")
	}
	fmt.Println("I think we connected");
}
