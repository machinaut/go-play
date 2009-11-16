// IRC - Internet Relay Chat package provides a convenient method of writing IRC clients and chat bots (RFC 2812).
package main

import (
    "./irc";
    "log"; // TODO: implement actual Loggers
    "time";
)


func main() {
    // Dial freenode
    server, err := irc.IRCDial("tcp","","irc.freenode.net:6667"); // TODO: make this not hardcoded
    if err != nil { log.Exit("Dialing error:", err); }
    // Login to the server
    server.Login("go_bot", irc.FlagInvisible);
    // Send a PM to NickServ to identify
    server.PrivMsg("NickServ","identify go_bot turing");
    // Join a chat
    bottest, _ := server.Join("#bottest"); // TODO: log the errors
    // Send the chat a message
    bottest.Write("hi guys!");
/*
    // Talk to server (loop forever)
    connReader := bufio.NewReader(conn);
    for i := 0; i < 100; i++ {
        response, err := connReader.ReadString('\n');
        if err != nil {
            log.Exit("Error reading from connection:", err)
        }
        log.Stdoutf("RECEIVED: %s", strings.TrimSpace(response));
        if response[0] != ':' { //not a private message
            wd := strings.Split(response, " ", 2);
            log.Stdout("Got Message ", wd[0]);
            switch wd[0] { // Message Type
            case "PING":
                // TODO: find a better way to remove leading character in string
                pongServer := string(strings.Bytes(wd[1])[1:len(wd[1])]);
                pong := "PONG " + pongServer + "\r\n";
                log.Stdout("SENT: ", pong);
                conn.Write(strings.Bytes(pong));
            }
        }
    }
    log.Stdout("Done reading response");
*/
    // We're done with the connection, close it
    time.Sleep(3e6);
    server.Close();
}
