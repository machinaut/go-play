// IRC - Internet Relay Chat package provides a convenient method of writing IRC clients and chat bots (RFC 2812).
// TODO: operate on byte[]s instead of strings for most stuff
package irc

import (
    "net";
    "os";
    "bufio";
    "strings";
    "strconv";
    "log";
)

// User modes for logging in, you probably just want Invisible on
// (and Wallops off)
type UserMode int

const (
    FlagWallops   = 1 << 2; // Get wallops (globally broadcast messages)
    FlagInvisible = 1 << 3; // User doesn't show up in WHO listings
)

// IRCConn is an (implementation of the net.Conn interface) for IRC network connections.
// Not there yet but working on it. TODO: actually fulfil the interface
type IRCConn struct {
    Nickname string;
    Usermode int;
    Conn     net.Conn;
}

// DialIRC() is like net.Dial() but can only connect to IRC networks
// and returns a IRCConn structure.
func IRCDial(netc, laddr, raddr string) (c *IRCConn, err os.Error) {
    log.Stdout("Dialing");
    conn, err := net.Dial(netc, laddr, raddr);
    if err != nil {
        return nil, err
    }
    (*c).Conn = conn;
    return c, nil;
}

// Write writes data to the TCP connection.
// TODO: Make this copy TCPConn's Write
func (c *IRCConn) Write(b []byte) (n int, err os.Error) {
    return c.Conn.Write(b)
}

// Read reads data from the TCP connection.
// TODO: Make this copy TCPConn's Read
func (c *IRCConn) Read(b []byte) (n int, err os.Error) {
    return c.Conn.Read(b)
}

// Read a line from the TCP connection.
// Returns the line as a "\r\n"-terminated string.
// TODO: clean this up like Read/Write
func (c *IRCConn) ReadLine() (string, os.Error) {
    rd := bufio.NewReader(c.Conn);
    return rd.ReadString('\n');
}

// Login() sends a login message with a given nickname
// and a given usermode.
func (c *IRCConn) Login(nick string, u UserMode) os.Error {
    // Form NICK message
    loginMessage := "NICK " + nick + "\r\n";
    // Form USER message
    username := "gobot";                   // TODO: make this a parameter
    realname := "Go Programming Language"; // TODO: make this a paramater
    loginMessage += "USER " + username + " " + strconv.Itoa(int(u)) + " * :" + realname + "\r\n";
    // Write to connection and return result
    _, err := c.Write(strings.Bytes(loginMessage)); // TODO: check actual count somehow
    return err;
}

// PrivMsg() sends a private message to the given user
func (c *IRCConn) PrivMsg(recipient, message string) os.Error {
    privMessage := "PRIVMSG " + recipient + " :" + message + "\r\n";
    // Write to connection and return result
    _, err := c.Write(strings.Bytes(privMessage)); // TODO: check bytes written (see net's TCP shit)
    return err;
}

// IRCChan - IRC Channel
type IRCChan struct {
    Conn *IRCConn; // IRC Connection this channel is on
    Chan string;   // IRC Channel name
}

// Join() - Join an IRC Channel on an IRC Connection
func (c *IRCConn) Join(channame string) (ch *IRCChan, err os.Error) {
    joinMessage := "JOIN " + channame + "\r\n";
    if _, err := c.Write(strings.Bytes(joinMessage)); err != nil {
        return nil, err
    }
    (*ch).Chan = channame;
    (*ch).Conn = c;
    return ch, err;
}

// Close() - close connection to server neatly
func (c *IRCConn) Close() os.Error { return c.Conn.Close() }

// Write() - write to an irc channel
// TODO: clean this up like the othre reads/writes
func (ch *IRCChan) Write(mesg string) os.Error {
    return ch.Conn.PrivMsg(ch.Chan, mesg)
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