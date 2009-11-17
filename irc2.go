package main

import (
    "net";
    "log";
    "os";
    "strings";
    "strconv";
    "time";
)

// IRC server connnection, has a simple logger and a TCP connection
type IRCConn struct {
    *net.Conn;
    *log.Logger;
}

// Constructor for IRCConn
// TODO: would it be cleaner to just oneline this in the code itself?
func NewIRCConn(conn *net.Conn, logger *log.Logger) *IRCConn {
    return &IRCConn{conn, logger}
}

// Send a message to the server
// Generic enough to be used to send any message
func (c *IRCConn) Mesg(mesgType string, params []string) (n int, err os.Error) {
    message := mesgType + " " + strings.Join(params, " ") + "\r\n";
    c.Log("SENT:", strings.TrimSpace(message));
    return c.Write(strings.Bytes(message));
}

// Send NICK message to server, securing a nickname
func (c *IRCConn) Nick(nick string) (n int, err os.Error) {
    params := make([]string, 1);
    params[0] = nick;
    return c.Mesg("NICK", params);
}

// Send USER message to server
// user : username requested
// name : real name
// mode : numeric(int) with 2 bit flags
//    'w'     bit 2 : recieve wallops
//    'i'     bit 3 : invisible
// (if in doubt use 8 as the mode)
func (c *IRCConn) User(user, name string, mode int) (n int, err os.Error) {
    params := make([]string, 4);
    params[0] = user;               // Username
    params[1] = strconv.Itoa(mode); // User Mode
    params[2] = "*";                // Unused parameter in the IRC protocol
    params[3] = ":" + name;         // Real Name
    return c.Mesg("USER", params);
}

// Login to an IRC server
// Uses the Nick and User methods to ident with some defaults
func (c *IRCConn) Login(nick string) (n int, err os.Error) {
    i, err := c.Nick(nick);
    if err != nil {
        return i, err
    }
    j, err := c.User("gobot", "Go Programming Language", 8);
    return i + j, err;
}

// Join an IRC Chatroom
// Chat names are usually of the form "#name" (e.g. "#fedora")
func (c *IRCConn) Join(name string) (n int, err os.Error) {
    params := make([]string, 1);
    params[0] = name;
    return c.Mesg("JOIN", params);
}

func main() {
    // Allocate a logger and a connection for our IRCConn
    logger := log.New(os.Stdout, nil, "", log.Lok|log.Ltime);
    conn, err := net.Dial("tcp", "", "irc.freenode.net:6667");
    if err != nil {
        log.Exit("Error dialing: ", err)
    }
    // Make new IRCConn from logger and tcp conn
    irc := NewIRCConn(&conn, logger);
    irc.Log("Hello, world!");

    // Login
    irc.Login("goo_bot");
    irc.Join("#ncsulug");

    time.Sleep(10e9); // 10 Seconds
    irc.Log("Closing");
    irc.Close();
}
