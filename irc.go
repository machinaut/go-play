// IRC - Internet Relay Chat package provides a convenient method of writing IRC clients and chat bots (RFC 2812).
package irc 

import (
    "net";
    "strings";
    "log"; // TODO: implement actual Loggers
    "bufio";
    "strconv";
)

// User modes for logging in, you probably just want Invisible on
// (and Wallops off)
type UserMode int

const (
    FlagWallops = 1 << 2; // Get wallops (globally broadcast messages)
    FlagInvisible = 1 << 3; // User doesn't show up in WHO listings
)

// IRCConn is an implementation of the net.Conn interface for IRC network connections.
type IRCConn struct{
    Nickname string;
    Usermode int;
    tcpConn *net.TCPConn; 
    }

// DialIRC() is like net.Dial() but can only connect to IRC networks
// and returns a IRCConn structure.
func IRCDial(net, laddr, raddr string) (c IRCConn, err os.Error) {
    var conn *TCPConn;
    if conn, err := net.Dial(net, laddr, raddr); err != nil {
        return nil, err;
    }
    var c *IRCConn;
    c.tcpConn = conn;
    return c;
}

// Login() sends a login message with a given nickname
// and a given usermode.
func (c *IRCConn) Login(nickname string, usermode int) os.Error {
    // Form NICK message
    loginMessage := "NICK " + nickname + "\r\n";
    // Form USER message
    username := "gobot";        // TODO: make this a parameter
    realname := "Go Programming Language"; // TODO: make this a paramater
    loginMessage += "USER " + username + " " + strconv.Itoa(usermode) + " * :" + realname + "\r\n";
    // Write to connection and return result
    _, err := c.Write(strings.Bytes(loginMessage);
    return err;
}

// PrivMsg() sends a private message to the given user
func (c *IRCConn) PrivMsg(recipient, message string) os.Error {
    privMessage := "PRIVMSG " + recipient + " :" + message + "\r\n";
    // Write to connection and return result
    _, err := c.Write(strings.Bytes(loginMessage);
    return err;
}

// IRCChan - IRC Channel
type IRCChan struct{
    Conn IRCConn; // IRC Connection this channel is on
    Chan string;  // IRC Channel name
}

// Join() - Join an IRC Channel on an IRC Connection
func (c *IRCConn) Join(channame string) *IRCChan, os.Error {
    joinMessage := "JOIN " + channame + "\r\n";
    if _, err := c.tcpConn.Write(strings.Bytes(loginMessage)); err != nil {
        return nil, err;
    }
    var ircChan *IRCChan;
    ircChan.Chan = channame;
    ircConn.Conn = c;
    return ircConn, err;
}

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

    // We're done with the connection, close it
    conn.Close();
}
*/
