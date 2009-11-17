// IRC - Internet Relay Chat package provides a convenient method of writing IRC clients and chat bots (RFC 2812).
package main

import (
    "./irc";
    "log"; // TODO: implement actual Loggers
    "time";
)


func main() {
    // Dial freenode
    log.Stdout("Dialing server");
    server, err := irc.IRCDial("tcp", "", "irc.freenode.net:6667"); // TODO: make this not hardcoded
    if err != nil {
        log.Exit("Dialing error:", err)
    }
    // Login to the server
    log.Stdout("Logging in to server");
    server.Login("goo_bot", irc.FlagInvisible);
    // Send a PM to NickServ to identify
    log.Stdout("Identifying to Nickserv");
    server.PrivMsg("NickServ", "identify go_bot turing");
    // Join a chat
    log.Stdout("Joining #bottest");
    bottest, _ := server.Join("#bottest"); // TODO: log the errors
    // Send the chat a message
    log.Stdout("Greeting #bottest");
    bottest.Write("hi guys!");

    // We're done with the connection, close it
    log.Stdout("Sleeping before closing");
    time.Sleep(1e10);
    log.Stdout("Closing");
    server.Close();
}
