// IRC - Internet Relay Chat package provides a convenient method of writing IRC clients and chat bots (RFC 2812).
package main

import (
    "net";
    "strings";
    "log"; // TODO: implement actual Loggers
    "bufio";
    "strconv";
)


func main() {
    // Connect to the server
    server := "irc.freenode.net:6667"; // TODO: make this not hardcoded
    // First resolve the name:port to an address
    addr, err := net.ResolveTCPAddr(server);
    if err != nil {
        log.Exitf("Error resolving server '%s': %s", server, err)
    }
    log.Stdoutf("Successfully resolved server '%s'", server);
    // Then Dial the address we just found to connect
    conn, err := net.DialTCP("tcp", nil, addr);
    if err != nil {
        log.Exitf("Error dialing server '%s': %s", server, err)
    }
    log.Stdoutf("Connected to server '%s'", server);

    // Formulate a message to login to server
    // Password
    password := "turing";
    loginMessage := "PASS " + password + "\r\n";
    // Nickname
    nickname := "go_bot"; // TODO: make this a parameter
    loginMessage += "NICK " + nickname + "\r\n";
    // Username
    username := "turing";        // TODO: make this a parameter
    realname := "Alonzo Church"; // TODO: make this a paramater
    // Usermode
    invisible := 1 << 3; // TODO: type user modes int const(whatever)
    usermode := invisible;
    // Send login message to server
    loginMessage += "USER " + username + " " + strconv.Itoa(usermode) + " * :" + realname + "\r\n";

    // Send a private message
    recipient := "ajray";
    message := "hi!";
    loginMessage += "PRIVMSG " + recipient + " :" + message + "\r\n";

    // Send a private message to NickServ identifying us
    recipient = "NickServ";
    message = "identify " + nickname + " " + password;
    loginMessage += "PRIVMSG " + recipient + " :" + message + "\r\n";

    // Tell it a channame
    channame := "#bottest"; // TODO: make this a parameter
    loginMessage += "JOIN " + channame + "\r\n";
    conn.Write(strings.Bytes(loginMessage));
    log.Stdoutf("SENT: %s", loginMessage);

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

    // We're done with the connection, close it
    conn.Close();
}
